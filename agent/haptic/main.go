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
	"fmt"

	"encoding/json"
	"os"
	"strings"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/dullgiulio/pingo" // for plugins
)

var (
	G_TwoStageActivation bool = true // TODO Legacy behaviour toggler, should disappear in the future

	G_Account AccountParams // Used to store account info used by the workflow procedures
	G_User    User          // Used to store account info used by the workflow procedures

	// Workflow procedures instances
	G_ProcRegisterProxyUser = &ProcRegisterProxyUser{}
	G_ProcCreateWinUser     = &ProcCreateWinUser{}

	g_Db Db

	g_PluginLdap     *pingo.Plugin
	g_PluginIaas     *pingo.Plugin
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
	case "listusers":
		users, err := ListUsers()
		if err != nil {
			nan.ExitError(err)
		} else {
			for _, u := range users {
				fmt.Printf("%v\n", u)
			}
		}

	case "registeruser":
		adapter.RegisterUser(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case "activateuser":
		adapter.ActivateUser(os.Args[2])
	case "deleteuser":
		adapter.DeleteUser(os.Args[2])
	case "changepassword":
		//TODO
	case "serve":
		RunServer()
	}
}

func InitPlugin(pluginName string, ppPlugin **pingo.Plugin) {
	var ok bool
	var pluginJsonParams nan.PluginParams

	pluginNameLowercase := strings.ToLower(pluginName)
	pluginRpcName := strings.ToUpper(pluginName[0:1]) + pluginName[1:len(pluginName)]

	if pluginJsonParams, ok = nan.Config().Plugins[pluginNameLowercase]; !ok {

		if pluginJsonParams, ok = nan.Config().Plugins[pluginRpcName]; !ok {
			LogError("Plugin %s doesn't have a parameters section in config.json !", pluginName)
			return
		}
	}

	pluginPath := fmt.Sprintf("%s/plugins/%s/%s", nan.Config().CommonBaseDir,
		pluginNameLowercase, pluginNameLowercase)

	*ppPlugin = pingo.NewPlugin("tcp", pluginPath)
	if *ppPlugin == nil {
		nan.ExitErrorf(0, "Failed to create plugin %s", pluginRpcName)
	}

	Log("Starting plugin %s", pluginRpcName)
	(*ppPlugin).Start()

	pluginParams, e := json.Marshal(pluginJsonParams)

	if e != nil {
		LogError("Failed to unmarshall %s plugin params", pluginName)
		ExitError(nan.ErrConfigError)
	}

	resp := ""

	if e := (*ppPlugin).Call(pluginRpcName+".Configure", string(pluginParams), &resp); e != nil {
		// TODO Clarify error and string output
		Log("Error while configuring plugin %s : %s", pluginRpcName, e)
		// TODO activate this line when all plugins have a Configure method
		// ExitError(nan.ErrPluginError)
	}

	Log("Start plugin %s : DONE", pluginRpcName)
}

func SetupPlugins() {
	Log("Num plugins referenced in config : %d", len(nan.Config().Plugins))
	InitPlugin("Iaas", &g_PluginIaas)
	InitPlugin("Ldap", &g_PluginLdap)
	InitPlugin("Owncloud", &g_PluginOwncloud)
}

func StopPlugins() {
	if g_PluginOwncloud != nil {
		g_PluginOwncloud.Stop()
	}

	if g_PluginLdap != nil {
		g_PluginLdap.Stop()
	}
}
