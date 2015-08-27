package main

import (
	"fmt"

	//nan "nanocloud.com/zeroinstall/lib/libnan"
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

	if !g_Db.GetRegisteredUsersInfo(&regUsersInfo) {
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

	if g_Db.GetActivatedUsersInfo(&activeTacUsersInfo) {

		for _, v := range activeTacUsersInfo {
			fmt.Println(v)
		}

	}

}
