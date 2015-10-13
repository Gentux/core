/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"

	"encoding/json"
	"os"
	"strings"

	nan "nanocloud.com/core/lib/libnan"

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
		if len(os.Args) != 6 {
			fmt.Println("Command registeruser expects 4 arguments.\nUsage: haptic registeruser firstname lastname email password\n")
		} else {
			adapter.RegisterUser(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
		}
	case "activateuser":
		if len(os.Args) != 3 {
			fmt.Println("Command activateuser expects 1 argument.\nUsage: haptic activateuser email\n")
		} else {
			nan.PrintErrorJson(adapter.ActivateUser(os.Args[2]))
			os.Exit(0)
		}

	case "deleteuser":
		if len(os.Args) != 3 {
			fmt.Println("Command deleteuser expects only one argument : the user email.\nUsage: haptic deleteuser email\n")
		} else {
			//adapter.DeleteUser(os.Args[2])

			var params AccountParams = AccountParams{
				Email: os.Args[2]}

			DeleteUser(params)
		}
	case "changeuserpassword":
		adapter.UpdateUserPassword(os.Args[2], os.Args[3])

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
		LogError("while configuring plugin %s : %s", pluginRpcName, e)
		// TODO activate this line when all plugins have a Configure method
		ExitError(nan.ErrPluginError)
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
