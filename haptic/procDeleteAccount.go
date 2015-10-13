/*
 * Nanocloud community -- transform any application into SaaS solution
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
	nan "nanocloud.com/core/lib/libnan"
)

// ========================================================================================================================

var ()

// ========================================================================================================================
// Procedure: DeleteUser
//
// Does:
// - Checks Params
// - [OPTIONAL] : free application specific resources + storage
// - DeleteUser (LDAP, AD)
// ========================================================================================================================
func DeleteUser(accountParams AccountParams) *nan.Err {

	Log("Starting procedure DeleteUser")

	ValidateDeleteUserParams(accountParams)

	var bRegistered, bActive bool
	var err *nan.Err

	if bActive, err = g_Db.IsUserActivated(accountParams.Email); err != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	}

	// TODO insert here freeing of application specific resources
	G_ProcCreateWinUser.Undo(accountParams)

	// TODO : do not exit too early here and allow for situations where an account may have been improperly created,
	// thus not visible in db anymore, but still with remaining files inside studio directory, that need to be all deleted

	if bRegistered, err = g_Db.IsUserRegistered(accountParams.Email); err != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	} else {
		G_ProcRegisterProxyUser.Undo()
	}

	if bActive && !bRegistered {
		LogError("Corrupt account cleaned up : was ACTIVE but didn't look NOT REGISTERED")
	} else if !bRegistered && !bActive {
		return LogErrorCode(ErrAccountDoesNotExist)
	}

	if user, err := g_Db.GetUser(accountParams.Email); err != nil {
		return LogErrorCode(err)
	} else {
		return g_Db.DeleteUser(user)
	}
}

// ========================================================================================================================
// ValidateUserParams
// ========================================================================================================================

// Procedure [NOSIDEEFFECT] Check preconditions for valid/compliant account creation parameters
func ValidateDeleteUserParams(accountParams AccountParams) {

	nan.Debug("Verifying parameters to delete %s account", accountParams.Email)

	if !nan.ValidEmail(accountParams.Email) {
		LogErrorCode(nan.ErrPbWithEmailFormat)
	}
}
