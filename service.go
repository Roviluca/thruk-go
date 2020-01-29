package thruk

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Service struct {
	FILE                       string   `json:":FILE"`
	ID                         string   `json:":ID,omitempty"`
	PEERKEY                    string   `json:":PEER_KEY,omitempty"`
	READONLY                   int      `json:":READONLY,omitempty"`
	TYPE                       string   `json:":TYPE"`
	ActionURL                  string   `json:"action_url,omitempty"`
	ActiveChecksEnabled        string   `json:"active_checks_enabled,omitempty"`
	CheckCommand               string   `json:"check_command,omitempty"`
	CheckFreshness             string   `json:"check_freshness,omitempty"`
	CheckInterval              string   `json:"check_interval,omitempty"`
	CheckPeriod                string   `json:"check_period,omitempty"`
	ContactGroups              []string `json:"contact_groups,omitempty"`
	Contacts                   []string `json:"contacts,omitempty"`
	DisplayName                string   `json:"display_name,omitempty"`
	EventHandler               []string `json:"event_handler,omitempty"`
	EventHandlerEnabled        string   `json:"event_handler_enabled,omitempty"`
	FirstNotificationDelay     string   `json:"first_notification_delay,omitempty"`
	FlapDetectionEnabled       string   `json:"flap_detection_enabled,omitempty"`
	FlapDetectionOptions       []string `json:"flap_detection_options,omitempty"`
	FreshnessThreshold         string   `json:"freshness_threshold,omitempty"`
	HighFlapThreshold          string   `json:"high_flap_threshold,omitempty"`
	HostName                   []string `json:"host_name,omitempty"`
	HostgroupName              []string `json:"hostgroup_name,omitempty"`
	IconImage                  string   `json:"icon_image,omitempty"`
	IconImageAlt               string   `json:"icon_image_alt,omitempty"`
	InitialState               string   `json:"initial_state,omitempty"`
	IsVolatile                 string   `json:"is_volatile,omitempty"`
	LowFlapThreshold           string   `json:"low_flap_threshold,omitempty"`
	MaxCheckAttempts           string   `json:"max_check_attempts,omitempty"`
	Name                       string   `json:"name"`
	Notes                      string   `json:"notes,omitempty"`
	NotesURL                   string   `json:"notes_url,omitempty"`
	NotificationInterval       string   `json:"notification_interval,omitempty"`
	NotificationOptions        []string `json:"notification_options,omitempty"`
	NotificationPeriod         string   `json:"notification_period,omitempty"`
	NotificationsEnabled       string   `json:"notifications_enabled,omitempty"`
	ObsessOverService          string   `json:"obsess_over_service,omitempty"`
	Parents                    []string `json:"parents,omitempty"`
	PassiveChecksEnabled       string   `json:"passive_checks_enabled,omitempty"`
	ProcessPerfData            string   `json:"process_perf_data,omitempty"`
	Register                   string   `json:"register,omitempty"`
	RetainNonstatusInformation string   `json:"retain_nonstatus_information,omitempty"`
	RetainStatusInformation    string   `json:"retain_status_information,omitempty"`
	RetryInterval              string   `json:"retry_interval,omitempty"`
	ServiceDescription         string   `json:"service_description,omitempty"`
	Servicegroups              []string `json:"servicegroups,omitempty"`
	StalkingOptions            []string `json:"stalking_options,omitempty"`
	Use                        []string `json:"use,omitempty"`
	WORKER                     string   `json:"_WORKER,omitempty"`
	FailurePredictionEnabled   string   `json:"failure_prediction_enabled,omitempty"`
}

func (t Thruk) GetService(id string) (Service, error) {
	var services []Service
	if id == "" {
		return Service{}, ErrorInvalidInput
	}
	resp, err := t.GetURL("/" + t.SiteName + "/thruk/r/config/objects?:TYPE=service&:ID=" + id)
	failOnError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&services)
	failOnError(err)
	if len(services) == 0 {
		return Service{}, ErrorObjectNotFound
	}

	return services[0], nil
}

func (t Thruk) CreateService(service Service) (string, error) {
	if service.FILE == "" || service.TYPE == "" {
		return "", ErrorNeedFileTypeHost
	}

	bodyBytes, _ := json.Marshal(service)
	body := bytes.NewReader(bodyBytes)
	resp, err := t.PostURL("/"+t.SiteName+"/thruk/r/config/objects/", body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.Status)
	}
	defer resp.Body.Close()

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

func (t Thruk) DeleteService(id string) error {
	URL := "/" + t.SiteName + "/thruk/r/config/objects/" + id
	err := t.DeleteURL(URL)
	if err != nil {
		return err
	}

	return nil
}
