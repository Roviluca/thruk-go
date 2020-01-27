package thruk

import (
	"encoding/json"
)

type Host struct {
	FILE                       string   `json:":FILE"`
	ID                         string   `json:":ID,omitempty"`
	PEERKEY                    string   `json:":PEER_KEY,omitempty"`
	READONLY                   int      `json:":READONLY,omitempty"`
	TYPE                       string   `json:":TYPE"`
	WORKER                     string   `json:"_WORKER,omitempty"`
	ActiveChecksEnabled        string   `json:"active_checks_enabled,omitempty"`
	Address                    string   `json:"address,omitempty"`
	CheckCommand               string   `json:"check_command,omitempty"`
	CheckInterval              string   `json:"check_interval,omitempty"`
	CheckPeriod                string   `json:"check_period,omitempty"`
	EventHandlerEnabled        string   `json:"event_handler_enabled,omitempty"`
	FlapDetectionEnabled       string   `json:"flap_detection_enabled,omitempty"`
	MaxCheckAttempts           string   `json:"max_check_attempts,omitempty"`
	Name                       string   `json:"name"`
	NotificationInterval       string   `json:"notification_interval,omitempty"`
	NotificationOptions        []string `json:"notification_options,omitempty"`
	NotificationPeriod         string   `json:"notification_period,omitempty"`
	NotificationsEnabled       string   `json:"notifications_enabled,omitempty"`
	ProcessPerfData            string   `json:"process_perf_data,omitempty"`
	Register                   string   `json:"register,omitempty"`
	RetainNonstatusInformation string   `json:"retain_nonstatus_information,omitempty"`
	RetainStatusInformation    string   `json:"retain_status_information,omitempty"`
	RetryInterval              string   `json:"retry_interval,omitempty"`
	ActionURL                  string   `json:"action_url,omitempty"`
	FailurePredictionEnabled   string   `json:"failure_prediction_enabled,omitempty"`
	Alias                      string   `json:"alias,omitempty"`
	Use                        []string `json:"use,omitempty"`
	TwoDCoords                 string   `json:"2d_coords,omitempty"`
	ThreeDCoords               string   `json:"3d_coords,omitempty"`
	CheckFreshness             string   `json:"check_freshness,omitempty"`
	ContactGroups              []string `json:"contact_groups,omitempty"`
	Contacts                   []string `json:"contacts,omitempty"`
	DisplayName                string   `json:"display_name,omitempty"`
	FirstNotificationDelay     string   `json:"first_notification_delay,omitempty"`
	FlapDetectionOptions       []string `json:"flap_detection_options,omitempty"`
	FreshnessThreshold         string   `json:"freshness_threshold,omitempty"`
	HighFlapThreshold          string   `json:"high_flap_threshold,omitempty"`
	HostName                   string   `json:"host_name,omitempty"`
	Hostgroups                 []string `json:"hostgroups,omitempty"`
	IconImage                  string   `json:"icon_image,omitempty"`
	IconImageAlt               string   `json:"icon_image_alt,omitempty"`
	InitialState               string   `json:"initial_state,omitempty"`
	LowFlapThreshold           string   `json:"low_flap_threshold,omitempty"`
	Notes                      string   `json:"notes,omitempty"`
	NotesURL                   string   `json:"notes_url,omitempty"`
	ObsessOverHost             string   `json:"obsess_over_host,omitempty"`
	Parents                    []string `json:"parents,omitempty"`
	PassiveChecksEnabled       string   `json:"passive_checks_enabled,omitempty"`
	StalkingOptions            []string `json:"stalking_options,omitempty"`
	StatusmapImage             string   `json:"statusmap_image,omitempty"`
	VrmlImage                  string   `json:"vrml_image,omitempty"`
}

func (t Thruk) GetHost(id string) (Host, error) {
	var hosts []Host
	if id == "" {
		return Host{}, ErrorInvalidInput
	}
	resp, err := t.GetURL("/" + t.SiteName + "/thruk/r/config/objects?:TYPE=host&:ID=" + id)
	failOnError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&hosts)
	failOnError(err)
	if len(hosts) == 0 {
		return Host{}, ErrorObjectNotFound
	}

	return hosts[0], nil
}
