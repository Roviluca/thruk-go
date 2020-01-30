package thruk

import (
	"bytes"
	"encoding/json"
	"errors"
)

var ErrorNeedFileTypeCommandName = errors.New("[ERROR] FILE, TYPE and CommandName must not be empty")

type Command struct {
	FILE        string `json:":FILE"`
	ID          string `json:":ID"`
	PEERKEY     string `json:":PEER_KEY"`
	READONLY    int    `json:":READONLY"`
	TYPE        string `json:":TYPE"`
	CommandLine string `json:"command_line"`
	CommandName string `json:"command_name"`
}

func (t Thruk) GetCommand(id string) (Command, error) {
	var commands []Command
	if id == "" {
		return Command{}, ErrorInvalidInput
	}
	resp, err := t.GetURL("/" + t.SiteName + "/thruk/r/config/objects?:TYPE=command&:ID=" + id)
	failOnError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&commands)
	failOnError(err)
	if len(commands) == 0 {
		return Command{}, ErrorObjectNotFound
	}

	return commands[0], nil
}

func (t Thruk) CreateCommand(command Command) (string, error) {
	if command.FILE == "" || command.TYPE == "" || command.CommandName == "" {
		return "", ErrorNeedFileTypeCommandName
	}

	bodyBytes, _ := json.Marshal(command)
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

//
func (t Thruk) DeleteCommand(id string) error {
	URL := "/" + t.SiteName + "/thruk/r/config/objects/" + id
	err := t.DeleteURL(URL)
	if err != nil {
		return err
	}

	return nil
}
