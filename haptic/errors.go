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
	ErrIssueWithDbResult        = nan.NewExitCode(16, "Could not parse database result")

	ErrMaxNumAccountsRegistered = nan.NewExitCode(4, "System has reached maximum number of registrations")
	ErrMaxNumAccountsReached    = nan.NewExitCode(4, "System has reached maximum number of accounts")

	OkAccountBeingCreated   = nan.NewExitCode(1, "The account is being created")
	OkAccountBeingDeleted   = nan.NewExitCode(1, "The account is being deleted")
	OkAccountBeingActivated = nan.NewExitCode(1, "The account is being activated")
)
