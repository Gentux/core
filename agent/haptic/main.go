package main

import (
	"fmt"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/dullgiulio/pingo" // for plugins

	"os"
)

var (
	G_TwoStageActivation bool = true // TODO Legacy behaviour toggler, should disappear in the future

	G_Account AccountParams // Used to store account info used by the workflow procedures

	// Workflow procedures instances
	G_ProcRegisterProxyUser = &ProcRegisterProxyUser{}
	G_ProcCreateWinUser     = &ProcCreateWinUser{}

	g_Db *Db

	g_PluginLdap     *pingo.Plugin
	g_PluginOwncloud *pingo.Plugin

	// Aliases useful for this package so that we don't have to have to prefix them with nan all the time
	ExitOk    = nan.ExitOk
	ExitError = nan.ExitError
	Log       = nan.Log
	LogError  = nan.LogError
)

func TestLdap() {

	// 	g_PluginLdap.Call("Ldap.Configure", params, &resp)

	var resp string

	params := `{ "userid" : "fred101@m.fr" }`
	g_PluginLdap.Call("Ldap.ForceDisableAccount", params, &resp)

	fmt.Println("Resp 0: ", resp)

	params = `{ "userid" : "fred101@m.fr", "password" : "Secr3TPass1942+" }`
	fmt.Println("Testing ldap with:", params)
	g_PluginLdap.Call("Ldap.AddUser", params, &resp)

	if resp[0] != '$' {
		LogError("Failed to add LDAP user, got output: <%s>. Retrying for user <%s> and password <%s>", resp, G_Account.Email, G_Account.Password)
	} else {
		fmt.Println("Resp 1: ", resp)
	}
}

// TODO remove hardcoded : intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10
func main() {

	// NOTE : libnan has an init() func that's already been called at this point and loaded the configuration file

	if len(os.Args) <= 1 {
		// TODO Print help message
		return
	}

	SetupPlugins()
	defer StopPlugins()

	Log("Initialise database driver")
	InitialiseDb()
	defer ShutdownDb()

	switch os.Args[1] {
	case "reg":
		adapter.RegisterUser(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case "activ":
		adapter.ActivateUser(os.Args[2])
	case "del":
		adapter.DeleteUser(os.Args[2])
	case "serve":
		RunServer()
	}
}

func SetupPlugins() {

	Log("Num plugins referenced in config : %d", len(nan.Config().Plugins))

	// LDAP
	g_PluginLdap = pingo.NewPlugin("tcp", "plugins/ldap/ldap")
	if g_PluginLdap == nil {
		nan.ExitErrorf(0, "Failed to start plugin Ldap")
	}

	// Owncloud
	g_PluginOwncloud = pingo.NewPlugin("tcp", "plugins/owncloud/owncloud")
	if g_PluginOwncloud == nil {
		nan.ExitErrorf(0, "Failed to start plugin Owncloud")
	}

	Log("Start plugin Ldap")
	g_PluginLdap.Start()
	Log("Start plugin Ldap : DONE")

	Log("Start plugin Owncloud")
	g_PluginOwncloud.Start()
	Log("Start plugin Owncloud : DONE")
}

func StopPlugins() {
	if g_PluginOwncloud != nil {
		g_PluginOwncloud.Stop()
	}

	if g_PluginLdap != nil {
		g_PluginLdap.Stop()
	}
}
