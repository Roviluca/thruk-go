package thruk

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Servicegroup struct {
	FILE                string   `json:":FILE"`
	ID                  string   `json:":ID,omitempty"`
	PEERKEY             string   `json:":PEER_KEY,omitempty"`
	READONLY            int      `json:":READONLY,omitempty"`
	TYPE                string   `json:":TYPE"`
	ActionURL           string   `json:"action_url,omitempty"`
	Alias               string   `json:"alias,omitempty"`
	Members             []string `json:"members,omitempty"`
	Name                string   `json:"name"`
	Notes               string   `json:"notes,omitempty"`
	NotesURL            string   `json:"notes_url,omitempty"`
	Register            string   `json:"register,omitempty"`
	ServicegroupMembers []string `json:"servicegroup_members,omitempty"`
	ServicegroupName    string   `json:"servicegroup_name,omitempty"`
	Use                 []string `json:"use,omitempty,omitempty"`
}

func (t Thruk) GetServicegroup(id string) (Servicegroup, error) {
	var servicegroups []Servicegroup
	if id == "" {
		return Servicegroup{}, ErrorInvalidInput
	}
	resp, err := t.GetURL("/" + t.SiteName + "/thruk/r/config/objects?:TYPE=servicegroup&:ID=" + id)
	failOnError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&servicegroups)
	failOnError(err)
	if len(servicegroups) == 0 {
		return Servicegroup{}, ErrorObjectNotFound
	}

	return servicegroups[0], nil
}

func (t Thruk) CreateServicegroup(servicegroup Servicegroup) (string, error) {
	if servicegroup.FILE == "" || servicegroup.TYPE == "" {
		return "", ErrorNeedFileTypeHost
	}

	bodyBytes, _ := json.Marshal(servicegroup)
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

func (t Thruk) DeleteServicegroup(id string) error {
	URL := "/" + t.SiteName + "/thruk/r/config/objects/" + id
	err := t.DeleteURL(URL)
	if err != nil {
		return err
	}

	return nil
}
