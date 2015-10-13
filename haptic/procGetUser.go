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

// ========================================================================================================================

const (
	ISO8601 = "2006-01-02T15:04:05Z0700"
)

var ()

func GetUserAccountRegistrationTime(_Email string) string {
	sUserParamsFilePath := fmt.Sprintf("%s/studio/%s/account_params", nan.Config().CommonBaseDir, _Email)

	sFileModTimeISO8601 := ""

	if fileInfo, err := os.Stat(sUserParamsFilePath); err != nil {
		LogError("Failed to access user registration data for email : %s", G_Account.Email)
	} else {
		sFileModTimeISO8601 = fileInfo.ModTime().Format(ISO8601)
	}

	return sFileModTimeISO8601
}

// ========================================================================================================================
// Procedure: CheckAccount
//
// Does:
// - Check Params
// - [OPTIONAL] checks status of application specific resources, eg. CheckConsulAgent()
// ========================================================================================================================

type CheckFullAccountParams struct {
	Email string
}

func CheckAccount(p CheckFullAccountParams) {

	G_Account.Email = p.Email

	InitialiseDb()
	defer ShutdownDb()

	ValidateCheckAccountParams()

	// Refuse check if user account not registered
	bRegistered, err := g_Db.IsUserRegistered(G_Account.Email)
	if err != nil {
		LogErrorCode(ErrIssueWithAccountsDb)
		return
	} else if !bRegistered {
		LogErrorCode(ErrAccountDoesNotExist)
		return
	}

	// Refuse check if user account not activated
	if active, err := g_Db.IsUserActivated(G_Account.Email); err != nil {
		LogErrorCode(ErrIssueWithAccountsDb)
		return

	} else if !active {

		// TODO insert here checking of application specific resources, eg. CheckConsulAgent()

		// Registered but not active : return descriptive message + registration time

		sFileModTimeISO8601 := GetUserAccountRegistrationTime(G_Account.Email)

		jsonStdout := fmt.Sprintf(`{"code":"4","message":"%s","timestamp":"%s"}`, ErrAccountNotActivated.Message, sFileModTimeISO8601)
		fmt.Println(jsonStdout)
		os.Exit(4)
	}

	jsonStdout := fmt.Sprintf(`{"code":"1","message":"Running"`)

	fmt.Println(jsonStdout)
}

func ValidateCheckAccountParams() {

	nan.Debug("Verifying parameters to check account for: %s", G_Account.Email)

	if !nan.ValidEmail(G_Account.Email) {
		LogErrorCode(nan.ErrPbWithEmailFormat)
	}
}
