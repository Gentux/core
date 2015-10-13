package main

import (
	"encoding/base64"
	"encoding/json"

	"github.com/hypersleep/easyssh"
)

// TODO
// ssh Administrator@10.20.30.40 -p 1119 "powershell.exe Get-RDRemoteApp"
// ssh Administrator@10.20.30.40 -p 1119 "powershell.exe -Command \"Get-RDRemoteApp | ConvertTo-Json\"" > /tmp/windowsShit.json
// https://technet.microsoft.com/fr-fr/library/jj215454.aspx

// ========================================================================================================================
// TYPES
// =====

type ApplicationParams struct {
	CollectionName string
	Alias          string
	DisplayName    string
	IconContents   []uint8
	FilePath       string
}

// ========================================================================================================================
// Procedure: ListUsers
//
// Does:
// - Return list of users
// ========================================================================================================================
func ListApplications() string {

	var (
		applications []ApplicationParams
		res          []byte
	)

	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:     "Administrator",
		Server:   "10.20.30.40",
		Port:     "1119",
		Password: "password",
	}

	// Call Run method with command you want to run on remote server.
	response, err := ssh.Run("powershell.exe -Command \"Get-RDRemoteApp | ConvertTo-Json -Compress\"")
	if err != nil {
		panic("Can't run remote command: " + err.Error())
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
	fmt.Println(Alias)
	var powershellCmd string

	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:     "Administrator",
		Server:   "10.20.30.40",
		Port:     "1119",
		Password: "password",
	}

	// TODO Parametrize this
	powershellCmd = fmt.Sprintf("powershell.exe -Command \"Remove-RDRemoteApp -Alias %s -CollectionName %s -Force\"", Alias, "winadapps")
	fmt.Println(powershellCmd)

	// Call Run method with command you want to run on remote server.
	_, err := ssh.Run(powershellCmd)
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	}
}
