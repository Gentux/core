package main

import (
	nan "nanocloud.com/zeroinstall/lib/libnan"
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
	} else {
		// TODO insert here freeing of application specific resources
		//G_ProcCreateWinUser.Undo(accountParams)
	}

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
