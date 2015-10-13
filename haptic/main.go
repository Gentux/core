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

	"os"

	nan "nanocloud.com/core/lib/libnan"
)

var (
	G_TwoStageActivation bool = true // TODO Legacy behaviour toggler, should disappear in the future

	G_Account AccountParams // Used to store account info used by the workflow procedures
	G_User    User          // Used to store account info used by the workflow procedures

	// Workflow procedures instances
	G_ProcRegisterProxyUser = &ProcRegisterProxyUser{}
	G_ProcCreateWinUser     = &ProcCreateWinUser{}

	g_Db Db

	g_PluginsMap PluginsMaps_t

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
