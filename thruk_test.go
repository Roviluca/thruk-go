package thruk

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"testing"
	"time"
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
	if err != nil {
		t.Error(err)
	}
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

		resp, err := thruk.GetURL(fmt.Sprintf("%s:%s", URL, ""))
		testCallError(t, err, resp.StatusCode)
	})
	t.Run("Can login with basic auth and list API root path", func(t *testing.T) {
		URL := startThrukContainer(t)
		skip_ssl_check := true
		thruk := newThruk(URL, omdTestUserName, omdTestPassword, skip_ssl_check)

		resp, err := thruk.GetURL(fmt.Sprintf("%s%s", URL, "/demo/thruk/r/"))
		testCallError(t, err, resp.StatusCode)
	})
}

type thruk struct {
	URL      string
	client   http.Client
	username string
	password string
}

func (t thruk) GetURL(URL string) (*http.Response, error) {
	req,err  := http.NewRequest("GET",t.URL,nil)
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}
	req.SetBasicAuth(t.username,t.password)
	resp, err := t.client.Do(req)
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}
	return resp, err
}
func newThruk(URL, username, password string, skipTLS bool) thruk {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
	}
	return thruk{
		URL: URL,
		client: http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
		username: username,
		password: password,
	}
}
