package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hypersleep/easyssh"

	nan "nanocloud.com/zeroinstall/lib/libnan"
)

type ApplicationParams struct {
	CollectionName string
	Alias          string
	DisplayName    string
	IconContents   []uint8
	FilePath       string
}

// ========================================================================================================================
// Procedure: ListApplications
//
// Does:
// - Return list of applications published by Active Directory
// ========================================================================================================================
func ListApplications() string {

	var (
		applications []ApplicationParams
		res          []byte
	)

	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:     nan.Config().AppServer.User,
		Server:   nan.Config().AppServer.Server,
		Port:     strconv.Itoa(nan.Config().AppServer.Port),
		Password: nan.Config().AppServer.Password,
	}

	// Call Run method with command you want to run on remote server.
	response, err := ssh.Run("powershell.exe -Command \"Get-RDRemoteApp | ConvertTo-Json -Compress\"")
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	} else if response == "" {
		response = "[]"
	}

	if []byte(response)[0] != []byte("[")[0] {
		response = fmt.Sprintf("[%s]", response)
	}
	json.Unmarshal([]byte(response), &applications)
	for _, application := range applications {
		application.IconContents = []byte(base64.StdEncoding.EncodeToString(application.IconContents))
	}

	res, _ = json.Marshal(applications)
	return string(res)
}

// ========================================================================================================================
// Procedure: UnpublishApplication
//
// Does:
// - Unpublish specified applications from ActiveDirectory
// ========================================================================================================================
func UnpublishApplication(Alias string) {
	var powershellCmd string

	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:     nan.Config().AppServer.User,
		Server:   nan.Config().AppServer.Server,
		Port:     strconv.Itoa(nan.Config().AppServer.Port),
		Password: nan.Config().AppServer.Password,
	}

	// TODO Parametrize this
	powershellCmd = fmt.Sprintf("powershell.exe -Command \"Remove-RDRemoteApp -Alias %s -CollectionName %s -Force\"", Alias, "winadapps")

	// Call Run method with command you want to run on remote server.
	_, err := ssh.Run(powershellCmd)
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	}
}
