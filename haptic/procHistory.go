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
	"encoding/json"
	"fmt"

	nan "nanocloud.com/core/lib/libnan"
)

type History struct {
	ConnectionId string
	StartDate    string
	EndDate      string
}

// ========================================================================================================================
// Procedure: GetHistory
//
// Does:
// - Return list of history for all users
// ========================================================================================================================
func GetHistory() []History {

	var (
		vms       string
		histories []History

		pPluginHistory *Plugin
		err            *nan.Err
	)

	pPluginHistory, err = GetPlugin("history")
	if err != nil || pPluginHistory == nil {
		return histories
	}

	e := pPluginHistory.Call("History.GetList", vms, &histories)
	if e != nil {
		fmt.Println("Error calling History Plugin: ", e)
		return nil
	}

	return histories
}

// ========================================================================================================================
// Procedure: GetHistoryForUser
//
// Does:
// - Get history for specific user
// ========================================================================================================================
func GetHistoryForUser(email string) []History {
	var (
		res       string
		histories []History

		pPluginHistory *Plugin
		err            *nan.Err
	)

	pPluginHistory, err = GetPlugin("history")
	if err != nil || pPluginHistory == nil {
		return histories
	}

	// TODO this can't work as I don't send email
	e := pPluginHistory.Call("History.GetList", nil, &res)
	if e != nil {
		LogError("when calling History Plugin: ", e)
	}

	if e = json.Unmarshal([]byte(res), &histories); e != nil {
		LogError("Cannot unmarshal output from History Plugin")
		return nil
	}

	return histories
}

// ========================================================================================================================
// Procedure: AddHistory
//
// Does:
// - Add history entry for a specific user
// ========================================================================================================================
func AddHistory(history History) bool {
	var res string

	jsonHistory, e := json.Marshal(history)
	if e != nil {
		fmt.Println(e)
		return false
	}

	var pPluginHistory *Plugin
	var err *nan.Err

	pPluginHistory, err = GetPlugin("history")
	if err != nil || pPluginHistory == nil {
		return false
	}

	if e = pPluginHistory.Call("History.Add", string(jsonHistory), &res); e != nil {
		LogError("when calling History Plugin: ", e)
	}

	if res == "true" {
		return true
	} else {
		return false
	}
}
