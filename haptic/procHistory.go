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
	)

	err := g_PluginHistory.Call("History.GetList", vms, &histories)
	if err != nil {
		fmt.Println("Error calling History Plugin: ", err)
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
	)

	// TODO this can't work as I don't send email
	err := g_PluginHistory.Call("History.GetList", nil, &res)
	if err != nil {
		fmt.Println("Error calling History Plugin: ", err)
	}

	e := json.Unmarshal([]byte(res), &histories)
	if e != nil {
		fmt.Println("Cannot unmarshal output from History Plugin")
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

	jsonHistory, err := json.Marshal(history)
	if err != nil {
		fmt.Println(err)
		return false
	}

	jsonHistory, err = json.Marshal(history)
	err = g_PluginHistory.Call("History.Add", string(jsonHistory), &res)
	if err != nil {
		fmt.Println("Error calling History Plugin: ", err)
	}

	if res == "true" {
		return true
	} else {
		return false
	}
}
