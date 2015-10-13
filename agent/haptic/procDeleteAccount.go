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
func DeleteUser(p AccountParams) {

	G_Account = p

	// InitialiseDb()
	// defer ShutdownDb()

	Log("Starting procedure DeleteUser")

	ValidateDeleteUserParams()

	var bRegistered, bActive bool
	var err error

	if bActive, err = g_Db.IsUserActivated(G_Account.Email); err != nil {
		ExitError(ErrIssueWithAccountsDb)
	} else {
		// TODO insert here freeing of application specific resources
		G_ProcCreateWinUser.Undo()
	}

	// TODO : do not exit too early here and allow for situations where an account may have been improperly created,
	// thus not visible in db anymore, but still with remaining files inside studio directory, that need to be all deleted

	if bRegistered, err = g_Db.IsUserRegistered(G_Account.Email); err != nil {
		ExitError(ErrIssueWithAccountsDb)
	} else {
		G_ProcRegisterProxyUser.Undo()
	}

	if bActive && !bRegistered {
		LogError("Corrupt account cleaned up : was ACTIVE but didn't look NOT REGISTERED")
	} else if !bRegistered && !bActive {
		ExitError(ErrAccountDoesNotExist)
	}

	ExitOk(OkAccountBeingDeleted)
}

// ========================================================================================================================
// ValidateUserParams
// ========================================================================================================================

// Procedure [NOSIDEEFFECT] Check preconditions for valid/compliant account creation parameters
func ValidateDeleteUserParams() {

	nan.Debug("Verifying parameters to delete %s account", G_Account.Email)

	if !nan.ValidEmail(G_Account.Email) {
		ExitError(nan.ErrPbWithEmailFormat)
	}
}
