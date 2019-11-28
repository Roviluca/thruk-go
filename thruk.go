package thruk

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type thruk struct {
	URL      string
	client   http.Client
	username string
	password string
}

type ConfigObject struct {
	FILE                        string   `json:":FILE"`
	ID                          string   `json:":ID"`
	PEERKEY                     string   `json:":PEER_KEY"`
	READONLY                    int      `json:":READONLY"`
	TYPE                        string   `json:":TYPE"`
	CommandLine                 string   `json:"command_line,omitempty"`
	CommandName                 string   `json:"command_name,omitempty"`
	ActionURL                   string   `json:"action_url,omitempty"`
	Name                        string   `json:"name,omitempty"`
	ProcessPerfData             string   `json:"process_perf_data,omitempty"`
	Register                    string   `json:"register,omitempty"`
	Alias                       string   `json:"alias,omitempty"`
	Friday                      string   `json:"friday,omitempty"`
	Monday                      string   `json:"monday,omitempty"`
	Saturday                    string   `json:"saturday,omitempty"`
	Sunday                      string   `json:"sunday,omitempty"`
	Thursday                    string   `json:"thursday,omitempty"`
	TimeperiodName              string   `json:"timeperiod_name,omitempty"`
	Tuesday                     string   `json:"tuesday,omitempty"`
	Wednesday                   string   `json:"wednesday,omitempty"`
	ActiveChecksEnabled         string   `json:"active_checks_enabled,omitempty"`
	CheckCommand                string   `json:"check_command,omitempty"`
	CheckInterval               string   `json:"check_interval,omitempty"`
	IsVolatile                  string   `json:"is_volatile,omitempty"`
	MaxCheckAttempts            string   `json:"max_check_attempts,omitempty"`
	PassiveChecksEnabled        string   `json:"passive_checks_enabled,omitempty"`
	RetryInterval               string   `json:"retry_interval,omitempty"`
	WORKER                      string   `json:"_WORKER,omitempty"`
	Address                     string   `json:"address,omitempty"`
	CheckPeriod                 string   `json:"check_period,omitempty"`
	EventHandlerEnabled         string   `json:"event_handler_enabled,omitempty"`
	FlapDetectionEnabled        string   `json:"flap_detection_enabled,omitempty"`
	NotificationInterval        string   `json:"notification_interval,omitempty"`
	NotificationOptions         []string `json:"notification_options,omitempty"`
	NotificationPeriod          string   `json:"notification_period,omitempty"`
	NotificationsEnabled        string   `json:"notifications_enabled,omitempty"`
	RetainNonstatusInformation  string   `json:"retain_nonstatus_information,omitempty"`
	RetainStatusInformation     string   `json:"retain_status_information,omitempty"`
	HostNotificationCommands    []string `json:"host_notification_commands,omitempty"`
	HostNotificationOptions     []string `json:"host_notification_options,omitempty"`
	HostNotificationPeriod      string   `json:"host_notification_period,omitempty"`
	ServiceNotificationCommands []string `json:"service_notification_commands,omitempty"`
	ServiceNotificationOptions  []string `json:"service_notification_options,omitempty"`
	ServiceNotificationPeriod   string   `json:"service_notification_period,omitempty"`
	CheckFreshness              string   `json:"check_freshness,omitempty"`
	ObsessOverService           string   `json:"obsess_over_service,omitempty"`
	FailurePredictionEnabled    string   `json:"failure_prediction_enabled,omitempty"`
}

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func (t thruk) GetURL(URL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", t.URL+URL, nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	req.SetBasicAuth(t.username, t.password)
	resp, err := t.client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return resp, err
}

func (t thruk) GetConfigObject(s string) (object ConfigObject, err error) {
	var configObjects []ConfigObject
	resp, err := t.GetURL("/demo/thruk/r/config/objects?:ID=341d4")
	defer resp.Body.Close()
	failOnError(err)

	err = json.NewDecoder(resp.Body).Decode(&configObjects)
	failOnError(err)

	return configObjects[0], nil
}

func failOnError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
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
