package main

import (
	nan "nanocloud.com/zeroinstall/lib/libnan"
)

var (
	ErrAccountDoesNotExist      = nan.NewExitCode(3, "The account does not exist")
	ErrAccountExists            = nan.NewExitCode(6, "The account already exists")
	ErrAccountNotRegistered     = nan.NewExitCode(7, "Account not registered yet")
	ErrAccountNotActivated      = nan.NewExitCode(7, "Account not activated yet")
	ErrAccountActivated         = nan.NewExitCode(8, "Account already activated")
	ErrFirstnameNonCompliant    = nan.NewExitCode(9, "Problem with the firstname format")
	ErrLastnameNonCompliant     = nan.NewExitCode(10, "Problem with the lastname format")
	ErrInvalidLicenseFile       = nan.NewExitCode(11, "Problem with the provided licence file")
	ErrIssueWithAdServer        = nan.NewExitCode(12, "Problem with the Active Directory server")
	ErrIssueWithTacProvisioning = nan.NewExitCode(13, "Problem during TAC server provisionning")
	ErrIncorrectNumberOfParams  = nan.NewExitCode(14, "Received incorrect number of parameters")
	ErrIssueWithAccountsDb      = nan.NewExitCode(15, "Could not access users database")

	ErrMaxNumAccountsRegistered = nan.NewExitCode(4, "System has reached maximum number of registrations")
	ErrMaxNumAccountsReached    = nan.NewExitCode(4, "System has reached maximum number of accounts")

	OkAccountBeingCreated   = nan.NewExitCode(1, "The account is being created")
	OkAccountBeingDeleted   = nan.NewExitCode(1, "The account is being deleted")
	OkAccountBeingActivated = nan.NewExitCode(1, "The account is being activated")
)
