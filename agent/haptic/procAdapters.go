package main

import (
	"fmt"
)

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
	var (
		user_id  int
		username string
		users    []string
	)

	rows, err := g_Db.Query("select user_id, username from guacamole_user")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user_id, &username)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, username)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return users, nil
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
