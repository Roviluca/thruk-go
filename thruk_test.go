package thruk

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/assert"
	"net/http"
	"testing"
)

const omdTestUserName string = "omdadmin"
const omdTestPassword string = "omdadmin"
const siteName string = "demo"

func startThrukContainer(t *testing.T) string {
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "test_container",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"443/tcp"},
		WaitingFor:   wait.ForLog("Thruk server is up and running"),
	}
	thrukContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	// Retrieve the container IP
	ip, err := thrukContainer.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	// Retrieve the port mapped to port 443
	port, err := thrukContainer.MappedPort(ctx, "443")
	if err != nil {
		t.Error(err)
	}
	URL := fmt.Sprintf("https://%s:%s", ip, port.Port())

	return URL
}
func startThrukServerAndGetClient(t *testing.T) *Thruk {
	t.Helper()
	URL := startThrukContainer(t)
	skipSslCheck := true
	thruk := NewThruk(URL, siteName, omdTestUserName, omdTestPassword, skipSslCheck)
	return thruk
}
func testCallError(t *testing.T, err error, respCode int) {
	t.Helper()
	assert.NilError(t, err)
	if respCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, respCode)
	}
}

func Test_container_up_and_running(t *testing.T) {
	t.Run("Test_container_up_and_running_with_authentication", func(t *testing.T) {
		URL := startThrukContainer(t)

		client := newClient()

		resp, err := client.Get(fmt.Sprintf("%s:%s", URL, ""))
		testCallError(t, err, resp.StatusCode)
	})
}

func Test_thruk_client(t *testing.T) {
	t.Run("Can contact thruk service successfully.", func(t *testing.T) {
		URL := startThrukContainer(t)
		skipSslCheck := true
		thruk := NewThruk(URL, siteName, "", "", skipSslCheck)

		resp, err := thruk.GetURL("")
		testCallError(t, err, resp.StatusCode)
	})
	t.Run("Can login with basic auth and list API root path", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		resp, err := thruk.GetURL("/demo/thruk/r/")
		testCallError(t, err, resp.StatusCode)
	})
}

func Test_thruk_client_GetConfigObject(t *testing.T) {
	t.Run("Can't get a config object of empty id and returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetConfigObject("")
		assert.Error(t, err, "[ERROR] invalid input")
		assert.DeepEqual(t, object, ConfigObject{})
	})
	t.Run("Can get a config object from id", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetConfigObject("341d4")
		assert.NilError(t, err)
		assert.DeepEqual(t, object, ConfigObject{
			FILE:           "/omd/sites/demo/etc/naemon/conf.d/thruk_templates.cfg:53",
			ID:             "341d4",
			PEERKEY:        "48f1c",
			READONLY:       0,
			TYPE:           "timeperiod",
			Alias:          "24 Hours A Day, 7 Days A Week",
			Friday:         "00:00-24:00",
			Monday:         "00:00-24:00",
			Saturday:       "00:00-24:00",
			Sunday:         "00:00-24:00",
			Thursday:       "00:00-24:00",
			TimeperiodName: "thruk_24x7",
			Tuesday:        "00:00-24:00",
			Wednesday:      "00:00-24:00",
		})

	})
}

func Test_thruk_client_CreateConfigObject(t *testing.T) {
	t.Run("returns error when FILE and TYPE are empty in object", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateConfigObject(ConfigObject{})
		assert.Error(t, err, "[ERROR] FILE and TYPE must not be empty")
	})
	t.Run("returns error if thruk returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateConfigObject(ConfigObject{
			FILE: "asd.asd",
			TYPE: "not_existent",
		})
		assert.Error(t, err, "object not created")
	})
	t.Run("returns nil error and ID", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateConfigObject(ConfigObject{
			FILE:    "test.cfg",
			TYPE:    "host",
			Name:    "localhost",
			Alias:   "localhost",
			Address: "127.0.0.1",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}
		createdObject, err := thruk.GetConfigObject(id)
		assert.NilError(t, err)

		assert.Equal(t, id, createdObject.ID)
	})
}

func Test_thruk_client_apply_operations(t *testing.T) {
	t.Run("discard configs from disk should drop unsaved changes", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateConfigObject(ConfigObject{
			FILE:    "test.cfg",
			TYPE:    "host",
			Name:    "localhost",
			Alias:   "localhost",
			Address: "127.0.0.1",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}

		createdObject, err := thruk.GetConfigObject(id)
		assert.NilError(t, err)
		assert.Equal(t, id, createdObject.ID)

		err = thruk.DiscardConfigs()
		assert.NilError(t, err)
		createdObject, err = thruk.GetConfigObject(id)
		assert.Error(t, err, "[ERROR] Config Object not found")
		assert.DeepEqual(t, createdObject, ConfigObject{})

	})
	t.Run("objects must persists after save function has been called", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateConfigObject(ConfigObject{
			FILE:    "test.cfg",
			TYPE:    "host",
			Name:    "localhost",
			Alias:   "localhost",
			Address: "127.0.0.1",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}

		err = thruk.SaveConfigs()
		assert.NilError(t, err)
		savedObject, err := thruk.GetConfigObject(id)
		assert.NilError(t, err)
		assert.DeepEqual(t, savedObject.ID, id)
	})
	t.Run("reload config function should return nil for success", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		err := thruk.ReloadConfigs()
		assert.NilError(t, err)
	})
	t.Run("CheckConfig function should return true if saved configuration is valid", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		ok := thruk.CheckConfig()
		assert.Assert(t, ok)

	})
	t.Run("CheckConfig function should return false if saved configuration is not valid", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		if _, err := thruk.CreateConfigObject(ConfigObject{TYPE: "host", FILE: "xxx.cfg"}); err != nil {
			t.Fatalf("Could not create object, error: %v", err)
		}
		if err := thruk.SaveConfigs(); err != nil {
			t.Fatalf("Could not save object, error: %v", err)
		}

		ok := thruk.CheckConfig()
		assert.Assert(t, !ok)

	})
}

func Test_thruk_client_DeleteConfigObject(t *testing.T) {
	t.Run("delete object must return nil error if object exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateConfigObject(ConfigObject{
			FILE: "test.cfg",
			TYPE: "host",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		err := thruk.DeleteConfigObject(id)
		assert.NilError(t, err)
	})
	t.Run("delete object must remove an object that exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateConfigObject(ConfigObject{
			FILE: "test.cfg",
			TYPE: "host",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		thruk.DeleteConfigObject(id)
		thruk.SaveConfigs()
		_, err := thruk.GetConfigObject(id)
		assert.Error(t, err, "[ERROR] Object not found")
	})
}
