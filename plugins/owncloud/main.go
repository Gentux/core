package main

import (
	"encoding/json"

	"github.com/dullgiulio/pingo"
	//"nanocloud.com/lib/connectors/owncloud"
)

// Create an object to be exported

type CreateUserParams struct {
	Username, Password string
}

type Owncloud struct{}

// func (p *Owncloud) Configure(_protocol, _hostname, _adminLogin, _password string) error {
// 	owncloud.Configure(_protocol, _hostname, _adminLogin, _password)
// 	return nil
// }

//func (p *Owncloud) CreateUser(_params []string, _outMsg *string) error {
func (p *Owncloud) CreateUser(jsonParams string, _outMsg *string) error {

	var params CreateUserParams

	if err := json.Unmarshal([]byte(jsonParams), &params); err != nil {
		*_outMsg = "ERROR: failed to unmarshal Owncloud.CreateUserParams : " + err.Error()
	} else {

		Configure("https", "esifront.nanocloud.com", "drive_admin", "secr3t")

		*_outMsg = params.Username
	}
	//owncloud.CreateUser(params.Username, params.Password).Error()

	return nil
}

func main() {
	plugin := &Owncloud{}

	pingo.Register(plugin)

	pingo.Run()
}
