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
func testCallError(t *testing.T, err error, respCode int) {
	t.Helper()
	testError(t, err, nil)
	if respCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, respCode)
	}
}

func testError(t *testing.T, got error, want error) {
	t.Helper()
	if got != want {
		fmt.Printf("Got error %s, wanted %s", got, want)
		t.Error(got)
	}
}
func Test_container_up_and_running(t *testing.T) {
	t.Run("Test_container_up_and_running_with_authentication", func(t *testing.T) {
		URL := startThrukContainer(t)

		client := newClient()

		resp, err := client.Get(fmt.Sprintf("%s:%s", URL, ""))
		testCallError(t, err, resp.StatusCode)
	})

	t.Run("Test_container_up_and_running2", func(t *testing.T) {
		URL := startThrukContainer(t)

		client := newClient()

		resp, err := client.Get(fmt.Sprintf("%s:%s", URL, ""))
		testCallError(t, err, resp.StatusCode)
	})
}

func Test_thruk_client(t *testing.T) {
	t.Run("Can contact thruk service successfully.", func(t *testing.T) {
		URL := startThrukContainer(t)
		skip_ssl_check := true
		thruk := newThruk(URL, "", "", skip_ssl_check)

		resp, err := thruk.GetURL("")
		testCallError(t, err, resp.StatusCode)
	})
	t.Run("Can login with basic auth and list API root path", func(t *testing.T) {
		URL := startThrukContainer(t)
		skip_ssl_check := true
		thruk := newThruk(URL, omdTestUserName, omdTestPassword, skip_ssl_check)

		resp, err := thruk.GetURL("/demo/thruk/r/")
		testCallError(t, err, resp.StatusCode)
	})
	t.Run("Can't get a config object of empty id and returns error", func(t *testing.T) {
		URL := startThrukContainer(t)
		skip_ssl_check := true
		thruk := newThruk(URL, omdTestUserName, omdTestPassword, skip_ssl_check)

		object, err := thruk.GetConfigObject("")
		testError(t, err, errorInvalidInput)
		assert.DeepEqual(t, object, ConfigObject{})
	})
	t.Run("Can get a config object from id", func(t *testing.T) {
		URL := startThrukContainer(t)
		skip_ssl_check := true
		thruk := newThruk(URL, omdTestUserName, omdTestPassword, skip_ssl_check)

		object, err := thruk.GetConfigObject("341d4")
		testError(t, err, nil)
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
