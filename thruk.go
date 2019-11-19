package thruk

import (
	"crypto/tls"
	"net/http"
)

func newClient()*http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}
