package main

import (
	"fmt"

	"github.com/boltdb/bolt"

	nan "nanocloud.com/zeroinstall/lib/libnan"
)

// We wrap sql.DB in a user struct to which we can add our own methods
type Db struct {
	*bolt.DB
}

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

func (o adapter_t) GetUsers() ([]User, *nan.Err) {
	return ListUsers(), nil
}

func (o adapter_t) UpdateUserEmail(_PrevEmail, _NewEmail string) error {

	fmt.Println("TODO UpdateUserEmail")

	return nil
}

func (o adapter_t) UpdateUserPassword(_Email, _Password string) error {

	UpdateUserPassword(_Email, _Password)

	return nil
}

func (o adapter_t) DeleteUser(_Email string) error {

	var params AccountParams = AccountParams{
		Email: _Email}

	DeleteUser(params)

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
