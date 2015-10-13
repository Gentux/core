package main

import (
	"encoding/json"

	"github.com/dullgiulio/pingo"
)

type Owncloud struct{}

type ConfigureParams struct {
	Protocol, Url, Login, Password string
}

func (p *Owncloud) Configure(_jsonParams string, _outMsg *string) error {

	var params ConfigureParams

	if err := json.Unmarshal([]byte(_jsonParams), &params); err != nil {
		*_outMsg = "ERROR: failed to unmarshal Owncloud.AddUserParams : " + err.Error()
	} else {
		*_outMsg = "1"
	}

	return OwncloudConfigure(params)
}

type AddUserParams struct {
	Username, Password string
}

func (p *Owncloud) AddUser(jsonParams string, _outMsg *string) error {
	var params AddUserParams

	if err := json.Unmarshal([]byte(jsonParams), &params); err != nil {
		*_outMsg = "ERROR: failed to unmarshal Owncloud.AddUserParams : " + err.Error()
	} else {
		*_outMsg = params.Username
	}

	return OwncloudAddUser(params.Username, params.Password)
}

type DeleteUserParams struct {
	Username string
}

func (p *Owncloud) DeleteUser(jsonParams string, _outMsg *string) error {

	var params DeleteUserParams

	if err := json.Unmarshal([]byte(jsonParams), &params); err != nil {
		*_outMsg = "ERROR: failed to unmarshal Owncloud.DeleteUserParams : " + err.Error()
	} else {
		*_outMsg = params.Username
	}
	//owncloud.CreateUser(params.Username, params.Password).Error()

	return OwncloudDeleteUser(params.Username)
}

func main() {
	plugin := &Owncloud{}

	pingo.Register(plugin)

	pingo.Run()
}
