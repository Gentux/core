package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/hypersleep/easyssh"

	nan "nanocloud.com/zeroinstall/lib/libnan"
)

type GuacamoleXMLConfigs struct {
	XMLName xml.Name             `xml:configs`
	Config  []GuacamoleXMLConfig `xml:"config"`
}

type GuacamoleXMLConfig struct {
	XMLName  xml.Name            `xml:config`
	Name     string              `xml:"name,attr"`
	Protocol string              `xml:"protocol,attr"`
	Params   []GuacamoleXMLParam `xml:"param"`
}

type GuacamoleXMLParam struct {
	ParamName  string `xml:"name,attr"`
	ParamValue string `xml:"value,attr"`
}

type Connection struct {
	Hostname       string `xml:"hostname"`
	Port           string `xml:"port"`
	Username       string `xml:"username"`
	Password       string `xml:"password"`
	RemoteApp      string `xml:"remote-app"`
	ConnectionName string
}

type ApplicationParams struct {
	CollectionName string
	Alias          string
	DisplayName    string
	IconContents   []uint8
	FilePath       string
}

// ========================================================================================================================
// Procedure: CreateConnections
//
// Does:
// - Create all connections in DB for a particular user in order to use all applications
// ========================================================================================================================
func CreateConnections() error {

	type configs GuacamoleXMLConfigs
	var (
		applications []ApplicationParams
		connections  configs
	)

	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		User:     nan.Config().Apps.AppServer.User,
		Server:   nan.Config().Apps.AppServer.Server,
		Port:     strconv.Itoa(nan.Config().Apps.AppServer.SSHPort),
		Password: nan.Config().Apps.AppServer.Password,
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

	users, _ := g_Db.GetUsers()
	for _, user := range users {
		for _, application := range applications {
			if application.Alias == "hapticPowershell" {
				continue
			}

			connections.Config = append(connections.Config, GuacamoleXMLConfig{
				Name:     fmt.Sprintf("%s_%s", application.Alias, user.Email),
				Protocol: "rdp",
				Params: []GuacamoleXMLParam{
					GuacamoleXMLParam{
						ParamName:  "hostname",
						ParamValue: nan.Config().Apps.AppServer.Server,
					},
					GuacamoleXMLParam{
						ParamName:  "port",
						ParamValue: strconv.Itoa(nan.Config().Apps.AppServer.RDPPort),
					},
					GuacamoleXMLParam{
						ParamName:  "username",
						ParamValue: fmt.Sprintf("%s@%s", user.Sam, nan.Config().Apps.AppServer.WindowsDomain),
					},
					GuacamoleXMLParam{
						ParamName:  "password",
						ParamValue: user.Password,
					},
					GuacamoleXMLParam{
						ParamName:  "remote-app",
						ParamValue: fmt.Sprintf("||%s", application.Alias),
					},
				},
			})
		}
	}

	connections.Config = append(connections.Config, GuacamoleXMLConfig{
		Name:     "hapticDesktop",
		Protocol: "rdp",
		Params: []GuacamoleXMLParam{
			GuacamoleXMLParam{
				ParamName:  "hostname",
				ParamValue: nan.Config().Apps.AppServer.Server,
			},
			GuacamoleXMLParam{
				ParamName:  "port",
				ParamValue: strconv.Itoa(nan.Config().Apps.AppServer.RDPPort),
			},
			GuacamoleXMLParam{
				ParamName:  "username",
				ParamValue: fmt.Sprintf("%s@%s", nan.Config().Apps.AppServer.User, nan.Config().Apps.AppServer.WindowsDomain),
			},
			GuacamoleXMLParam{
				ParamName:  "password",
				ParamValue: nan.Config().Apps.AppServer.Password,
			},
		},
	})
	connections.Config = append(connections.Config, GuacamoleXMLConfig{
		Name:     "hapticPowershell",
		Protocol: "rdp",
		Params: []GuacamoleXMLParam{
			GuacamoleXMLParam{
				ParamName:  "hostname",
				ParamValue: nan.Config().Apps.AppServer.Server,
			},
			GuacamoleXMLParam{
				ParamName:  "port",
				ParamValue: strconv.Itoa(nan.Config().Apps.AppServer.RDPPort),
			},
			GuacamoleXMLParam{
				ParamName:  "username",
				ParamValue: fmt.Sprintf("%s@%s", nan.Config().Apps.AppServer.User, nan.Config().Apps.AppServer.WindowsDomain),
			},
			GuacamoleXMLParam{
				ParamName:  "password",
				ParamValue: nan.Config().Apps.AppServer.Password,
			},
			GuacamoleXMLParam{
				ParamName:  "remote-app",
				ParamValue: "||hapticPowershell",
			},
		},
	})

	output, err := xml.MarshalIndent(connections, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	if err = ioutil.WriteFile(nan.Config().Apps.AppServer.XMLConfigurationFile, output, 0777); err != nil {
		LogError("Failed to save connections in %s params : %s", nan.Config().Apps.AppServer.XMLConfigurationFile, err)
		return err
	}

	return nil
}

// ========================================================================================================================
// Procedure: ListApplications
//
// Does:
// - Return list of applications published by Active Directory
// ========================================================================================================================
func ListApplications() []Connection {

	var (
		guacamoleConfigs GuacamoleXMLConfigs
		connections      []Connection
		bytesRead        []byte
		err              error
	)

	go CreateConnections()

	if bytesRead, err = ioutil.ReadFile(nan.Config().Apps.AppServer.XMLConfigurationFile); err != nil {
		LogError("Failed to read connections params : %s", err)
	}

	err = xml.Unmarshal(bytesRead, &guacamoleConfigs)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	for _, config := range guacamoleConfigs.Config {
		var connection Connection

		for _, param := range config.Params {
			switch true {
			case param.ParamName == "hostname":
				connection.Hostname = param.ParamValue
			case param.ParamName == "port":
				connection.Port = param.ParamValue
			case param.ParamName == "username":
				connection.Username = param.ParamValue
			case param.ParamName == "password":
				connection.Password = param.ParamValue
			case param.ParamName == "remote-app":
				connection.RemoteApp = param.ParamValue
			}
		}
		connection.ConnectionName = config.Name

		if connection.RemoteApp == "" || connection.RemoteApp == "||hapticPowershell" {
			continue
		}

		connections = append(connections, connection)
	}

	return connections
}

// ========================================================================================================================
// Procedure: ListApplicationsForSamAccount
//
// Does:
// - Return list of applications available for a particular SAM account
// ========================================================================================================================
func ListApplicationsForSamAccount(sam string) []Connection {

	var (
		guacamoleConfigs GuacamoleXMLConfigs
		connections      []Connection
		bytesRead        []byte
		err              error
	)

	if bytesRead, err = ioutil.ReadFile(nan.Config().Apps.AppServer.XMLConfigurationFile); err != nil {
		LogError("Failed to read connections params : %s", err)
	}

	err = xml.Unmarshal(bytesRead, &guacamoleConfigs)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	for _, config := range guacamoleConfigs.Config {
		var connection Connection

		if connection.ConnectionName == "hapticPowershell" {
			continue
		}

		connection.ConnectionName = config.Name
		for _, param := range config.Params {
			switch true {
			case param.ParamName == "hostname":
				connection.Hostname = param.ParamValue
			case param.ParamName == "port":
				connection.Port = param.ParamValue
			case param.ParamName == "username":
				connection.Username = param.ParamValue
			case param.ParamName == "password":
				connection.Password = param.ParamValue
			case param.ParamName == "remote-app":
				connection.RemoteApp = param.ParamValue
			}
		}

		if connection.Username == fmt.Sprintf("%s@%s", sam, nan.Config().Apps.AppServer.WindowsDomain) {
			connections = append(connections, connection)
		}
	}

	return connections
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
		User:     nan.Config().Apps.AppServer.User,
		Server:   nan.Config().Apps.AppServer.Server,
		Port:     strconv.Itoa(nan.Config().Apps.AppServer.SSHPort),
		Password: nan.Config().Apps.AppServer.Password,
	}

	// TODO Parametrize this
	powershellCmd = fmt.Sprintf("powershell.exe -Command \"Remove-RDRemoteApp -Alias %s -CollectionName %s -Force\"", Alias, "winadapps")

	// Call Run method with command you want to run on remote server.
	_, err := ssh.Run(powershellCmd)
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	}
}

// ========================================================================================================================
// Procedure: SyncUploadedFile
//
// Does:
// - Upload user files to windows VM
// ========================================================================================================================
func SyncUploadedFile(Filename string) {

	ssh := &easyssh.MakeConfig{
		User:     nan.Config().Apps.AppServer.User,
		Server:   nan.Config().Apps.AppServer.Server,
		Port:     strconv.Itoa(nan.Config().Apps.AppServer.SSHPort),
		Password: nan.Config().Apps.AppServer.Password,
	}

	fmt.Println("SyncUploadedFile")

	// Call Scp method with file you want to upload to remote server.
	err := ssh.Scp(Filename)

	// Handle errors
	if err != nil {
		LogError("Can't run remote command: " + err.Error())
	} else {
		fmt.Printf("SCP upload success for file %s\n", Filename)
	}
}
