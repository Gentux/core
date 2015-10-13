package main

import (
	"fmt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// We wrap sql.DB in a user struct to which we can add our own methods
type Db struct {
	*sql.DB
}

type adapter_t struct{}

var (
	adapter adapter_t
)

func (o adapter_t) RegisterUser(_Firstname, _Lastname, _Email, _Password string) error {

	var params AccountParams = AccountParams{

		FirstName: _Firstname,
		LastName:  _Lastname,
		Email:     _Email,
		Password:  _Password}

	RegisterUser(params)

	return nil
}

func (o adapter_t) ActivateUser(_Email string) error {

	var params AccountParams = AccountParams{
		Email: _Email}

	ActivateUser(params)

	return nil
}

func (o adapter_t) GetUsers() ([]string, error) {

	return ListUsers(), nil
}

func (o adapter_t) UpdateUserEmail(_PrevEmail, _NewEmail string) error {

	fmt.Println("TODO UpdateUserEmail")

	return nil
}

func (o adapter_t) UpdateUserPassword(_Email, _Password string) error {

	fmt.Println("TODO UpdateUserPassword")

	return nil
}

func (o adapter_t) DeleteUser(_Email string) error {

	var params AccountParams = AccountParams{
		Email: _Email}

	DeleteUser(params)

	return nil
}

func (o adapter_t) GetApplications() (string, error) {

	return ListApplications(), nil
}

func (o adapter_t) UnpublishApp(Alias string) error {

	UnpublishApplication(Alias)

	return nil
}
