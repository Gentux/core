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
	nan "nanocloud.com/core/lib/libnan"
)

type adapter_t struct{}

var (
	adapter adapter_t
)

func (o adapter_t) RegisterUser(_Firstname, _Lastname, _Email, _Password string) *nan.Err {

	var params AccountParams = AccountParams{

		FirstName: _Firstname,
		LastName:  _Lastname,
		Email:     _Email,
		Password:  _Password}

	return RegisterUser(params)
}

func (o adapter_t) ActivateUser(_Email string) *nan.Err {

	var params AccountParams = AccountParams{
		Email: _Email}

	return ActivateUser(params)
}

func (o adapter_t) UpdateUserPassword(_Email, _Password string) *nan.Err {

	if UpdateUserPassword(_Email, _Password) != true {
		return nan.ErrPasswordNotUpdated
	}

	return nil
}

func (o adapter_t) GetApplications() ([]Connection, error) {

	return ListApplications(), nil
}

func (o adapter_t) GetApplicationsForSamAccount(sam string) ([]Connection, error) {

	return ListApplicationsForSamAccount(sam), nil
}

func (o adapter_t) UnpublishApp(Alias string) error {

	UnpublishApplication(Alias)

	return nil
}

func (o adapter_t) GetVmList() (string, error) {
	return ListVMs(), nil
}

func (o adapter_t) DownloadWindowsVm() (bool, error) {
	return DownloadWindowsVm(), nil
}

func (o adapter_t) DownloadStatus() (bool, error) {
	return DownloadStatus(), nil
}

func (o adapter_t) StartVm(vmName string) (bool, error) {
	return StartVm(vmName), nil
}
