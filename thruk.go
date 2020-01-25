package thruk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var ErrorInvalidInput = errors.New("[ERROR] invalid input")
var ErrorNeedFileAndType = errors.New("[ERROR] FILE and TYPE must not be empty")
var ErrorObjectNotFound = errors.New("[ERROR] Config Object not found")

type Thruk struct {
	URL      string
	client   http.Client
	username string
	password string
	SiteName string
}

type thrukResponse struct {
	Count   int            `json:"count"`
	Message string         `json:"message"`
	Objects []ConfigObject `json:"objects"`
}
type ConfigObject struct {
	FILE                        string   `json:":FILE"`
	ID                          string   `json:":ID,omitempty"`
	PEERKEY                     string   `json:":PEER_KEY,omitempty"`
	READONLY                    int      `json:":READONLY,omitempty"`
	TYPE                        string   `json:":TYPE"`
	CommandLine                 string   `json:"command_line,omitempty"`
	CommandName                 string   `json:"command_name,omitempty"`
	ActiveChecksEnabled         string   `json:"active_checks_enabled,omitempty"`
	CheckFreshness              string   `json:"check_freshness,omitempty"`
	CheckInterval               string   `json:"check_interval,omitempty"`
	CheckPeriod                 string   `json:"check_period,omitempty"`
	EventHandlerEnabled         string   `json:"event_handler_enabled,omitempty"`
	FailurePredictionEnabled    string   `json:"failure_prediction_enabled,omitempty"`
	FlapDetectionEnabled        string   `json:"flap_detection_enabled,omitempty"`
	IsVolatile                  string   `json:"is_volatile,omitempty"`
	MaxCheckAttempts            string   `json:"max_check_attempts,omitempty"`
	Name                        string   `json:"name,omitempty"`
	NotificationInterval        string   `json:"notification_interval,omitempty"`
	NotificationOptions         []string `json:"notification_options,omitempty"`
	NotificationPeriod          string   `json:"notification_period,omitempty"`
	NotificationsEnabled        string   `json:"notifications_enabled,omitempty"`
	ObsessOverService           string   `json:"obsess_over_service,omitempty"`
	PassiveChecksEnabled        string   `json:"passive_checks_enabled,omitempty"`
	ProcessPerfData             string   `json:"process_perf_data,omitempty"`
	Register                    string   `json:"register,omitempty"`
	RetainNonstatusInformation  string   `json:"retain_nonstatus_information,omitempty"`
	RetainStatusInformation     string   `json:"retain_status_information,omitempty"`
	RetryInterval               string   `json:"retry_interval,omitempty"`
	Alias                       string   `json:"alias,omitempty"`
	TimeperiodName              string   `json:"timeperiod_name,omitempty"`
	CheckCommand                string   `json:"check_command,omitempty"`
	Friday                      string   `json:"friday,omitempty"`
	Monday                      string   `json:"monday,omitempty"`
	Saturday                    string   `json:"saturday,omitempty"`
	Sunday                      string   `json:"sunday,omitempty"`
	Thursday                    string   `json:"thursday,omitempty"`
	Tuesday                     string   `json:"tuesday,omitempty"`
	Wednesday                   string   `json:"wednesday,omitempty"`
	ActionURL                   string   `json:"action_url,omitempty"`
	TwoDCoords                  string   `json:"2d_coords,omitempty"`
	ThreeDCoords                string   `json:"3d_coords,omitempty"`
	Address                     string   `json:"address,omitempty"`
	ContactGroups               []string `json:"contact_groups,omitempty"`
	Contacts                    []string `json:"contacts,omitempty"`
	DisplayName                 string   `json:"display_name,omitempty"`
	EventHandler                []string `json:"event_handler,omitempty"`
	FirstNotificationDelay      string   `json:"first_notification_delay,omitempty"`
	FlapDetectionOptions        []string `json:"flap_detection_options,omitempty"`
	FreshnessThreshold          string   `json:"freshness_threshold,omitempty"`
	HighFlapThreshold           string   `json:"high_flap_threshold,omitempty"`
	HostName                    string   `json:"host_name,omitempty"`
	Hostgroups                  []string `json:"hostgroups,omitempty"`
	IconImage                   string   `json:"icon_image,omitempty"`
	IconImageAlt                string   `json:"icon_image_alt,omitempty"`
	InitialState                string   `json:"initial_state,omitempty"`
	LowFlapThreshold            string   `json:"low_flap_threshold,omitempty"`
	Notes                       string   `json:"notes,omitempty"`
	NotesURL                    string   `json:"notes_url,omitempty"`
	ObsessOverHost              string   `json:"obsess_over_host,omitempty"`
	Parents                     []string `json:"parents,omitempty"`
	StalkingOptions             []string `json:"stalking_options,omitempty"`
	StatusmapImage              string   `json:"statusmap_image,omitempty"`
	Use                         []string `json:"use,omitempty"`
	VrmlImage                   string   `json:"vrml_image,omitempty"`
	WORKER                      string   `json:"_WORKER,omitempty"`
	HostNotificationCommands    []string `json:"host_notification_commands,omitempty"`
	HostNotificationOptions     []string `json:"host_notification_options,omitempty"`
	HostNotificationPeriod      string   `json:"host_notification_period,omitempty"`
	ServiceNotificationCommands []string `json:"service_notification_commands,omitempty"`
	ServiceNotificationOptions  []string `json:"service_notification_options,omitempty"`
	ServiceNotificationPeriod   string   `json:"service_notification_period,omitempty"`
	ServicegroupName            string   `json:"servicegroup_name,omitempty"`
	Servicegroups               []string `json:"servicegroups,omitempty"`
}

type reloadResponse []struct {
	Failed  bool   `json:"failed"`
	Output  string `json:"output"`
	PeerKey string `json:"peer_key"`
}

type deleteResponse struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

type checkResponse []struct {
	Failed bool   `json:"failed"`
	Output string `json:"output"`
}

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func (t Thruk) GetURL(URL string) (*http.Response, error) {
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

func (t Thruk) PostURL(URL string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", t.URL+URL, body)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil, err
	}
	req.SetBasicAuth(t.username, t.password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return nil, err
	}

	return resp, err
}

func (t Thruk) GetConfigObject(id string) (object ConfigObject, err error) {
	var configObjects []ConfigObject
	if id == "" {
		return object, ErrorInvalidInput
	}
	resp, err := t.GetURL("/" + t.SiteName + "/thruk/r/config/objects?:ID=" + id)
	defer resp.Body.Close()
	failOnError(err)

	err = json.NewDecoder(resp.Body).Decode(&configObjects)
	failOnError(err)
	if len(configObjects) == 0 {
		return ConfigObject{}, ErrorObjectNotFound
	}

	return configObjects[0], nil
}

func (t Thruk) CreateConfigObject(object ConfigObject) (id string, err error) {
	if object.FILE == "" || object.TYPE == "" {
		return "", ErrorNeedFileAndType
	}

	bodyBytes, _ := json.Marshal(object)
	body := bytes.NewReader(bodyBytes)
	resp, err := t.PostURL("/"+t.SiteName+"/thruk/r/config/objects/", body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.Status)
	}

	thrukResp := thrukResponse{}
	err = json.NewDecoder(resp.Body).Decode(&thrukResp)
	if err != nil {
		return "", err
	}
	if len(thrukResp.Objects) == 0 {
		return "", errors.New("object not created")
	}
	return thrukResp.Objects[0].ID, err
}

func (t Thruk) DiscardConfigs() error {
	resp, err := t.PostURL("/"+t.SiteName+"/thruk/r/config/discard", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http response code %d", resp.StatusCode)
	}
	return nil
}

func (t Thruk) SaveConfigs() error {
	resp, err := t.PostURL("/"+t.SiteName+"/thruk/r/config/save", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http response code %d", resp.StatusCode)
	}
	return nil
}

func (t Thruk) ReloadConfigs() error {
	reloadResp := reloadResponse{}
	resp, err := t.PostURL("/"+t.SiteName+"/thruk/r/config/reload", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http response code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&reloadResp)
	if err != nil {
		return err
	}
	if reloadResp[0].Failed {
		return fmt.Errorf("reload failed %v", reloadResp)
	}
	return nil
}

func (t Thruk) DeleteConfigObject(id string) error {
	URL := "/" + t.SiteName + "/thruk/r/config/objects/" + id
	err := t.DeleteURL(URL)
	if err != nil {
		return err
	}

	return nil
}

func (t Thruk) DeleteURL(URL string) error {
	req, err := http.NewRequest("DELETE", t.URL+URL, nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return err
	}
	req.SetBasicAuth(t.username, t.password)
	resp, err := t.client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}
	return nil
}

func (t Thruk) CheckConfig() bool {
	checkResult := checkResponse{}

	resp, err := t.PostURL("/"+t.SiteName+"/thruk/r/config/check", nil)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	json.NewDecoder(resp.Body).Decode(&checkResult)
	return !checkResult[0].Failed
}

func failOnError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
func NewThruk(URL, SiteName, username, password string, skipTLS bool) *Thruk {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
	}
	return &Thruk{
		URL:      URL,
		SiteName: SiteName,
		client: http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
		username: username,
		password: password,
	}
}
