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

	//nan "nanocloud.com/core/lib/libnan"
)

// ========================================================================================================================
// TYPES

type RegisteredUserInfo struct {
	Email        string `json:"user_id"`
	CreationTime string `json:"created_date"`
	Activated    bool   `json:"activated"`
}

type ActiveTacUserInfo struct {
	TacId        string `json:"user_id"`
	TacUrl       string `json:"tac_url"`
	CreationTime string `json:"created_date"`
}

// ========================================================================================================================

const ()

var ()

// ========================================================================================================================
// Procedure: listRegisteredUsers
//
// Does:
// -
// ========================================================================================================================
func ListRegisteredUsers() {

	// Retirer restriction sur non activated => inclure les utilisateurs actifs
	// TODO Ajouter active : false/true

	InitialiseDb()
	defer ShutdownDb()

	var regUsersInfo []RegisteredUserInfo

	if err := g_Db.GetRegisteredUsersInfo(&regUsersInfo); err != nil {
		LogErrorCode(err)
		return
	}

	sResult := `{"fields":["user_id","registration_date","activated"],"data":[`
	for idx, regInfo := range regUsersInfo {

		sResult += fmt.Sprintf(`["%s","%s", "%t"]`, regInfo.Email, regInfo.CreationTime, regInfo.Activated)

		if idx < len(regUsersInfo)-1 {
			sResult += ","
		}
	}
	sResult += `]}`

	fmt.Println(sResult)

}

// ========================================================================================================================
// Procedure: listRegisteredUsers
//
// Does:
// -
// ========================================================================================================================
func ListActiveUsers() {

	InitialiseDb()
	defer ShutdownDb()

	var activeTacUsersInfo []ActiveTacUserInfo

	if err := g_Db.GetActivatedUsersInfo(&activeTacUsersInfo); err != nil {
		LogErrorCode(err)
		return
	}

	for _, v := range activeTacUsersInfo {
		fmt.Println(v)
	}

}
