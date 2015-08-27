// Copyright 2015 Nanocloud SAS (Paris, France).
//
// Licensed under the XXX License, Version N (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http:// TODO
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// TODO Insert program overview
//

package main

import (
	"encoding/json"
	"os"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/dullgiulio/pingo" // for plugins
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
	ExitOk       = nan.ExitOk
	ExitError    = nan.ExitError
	Log          = nan.Log
	LogError     = nan.LogError
	LogErrorCode = nan.LogErrorCode
)

// TODO remove hardcoded : intra.nanocloud.com/Administrator%password", "//10.20.12.10
func main() {

	// NOTE : libnan has an init() func that's already been called at this point and loaded the configuration file

	if len(os.Args) <= 1 {
		return
	}

	SetupPlugins()
	defer StopPlugins()

	Log("Initialise database driver")
	InitialiseDb()
	defer ShutdownDb()

	switch os.Args[1] {
	case "registeruser":
		adapter.RegisterUser(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case "activateuser":
		nan.PrintErrorJson(adapter.ActivateUser(os.Args[2]))
		os.Exit(0)
	case "deleteuser":
		adapter.DeleteUser(os.Args[2])
	case "changepassword":
		//TODO
	case "serve":
		RunServer()
	}
}

func SetupPlugins() {
	Log("Num plugins referenced in config : %d", len(nan.Config().Plugins))

	// LDAP
	g_PluginLdap = pingo.NewPlugin("tcp", nan.Config().CommonBaseDir+"/plugins/ldap/ldap")
	if g_PluginLdap == nil {
		nan.LogError("Failed to start plugin Ldap")
	}

	// Owncloud
	g_PluginOwncloud = pingo.NewPlugin("tcp", nan.Config().CommonBaseDir+"/plugins/owncloud/owncloud")
	if g_PluginOwncloud == nil {
		nan.LogError("Failed to start plugin Owncloud")
	}

	Log("Start plugin Ldap")
	g_PluginLdap.Start()
	var (
		resp                 string
		pluginLdapJsonParams []byte
	)
	pluginLdapJsonParams, _ = json.Marshal(nan.Config().Plugins["Ldap"])
	err := g_PluginLdap.Call("Ldap.Configure", string(pluginLdapJsonParams), &resp)
	if err != nil {
		// TODO Clarify error and string output
		Log("Error while configuring plugin Ldap : %s", err)
	}
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
