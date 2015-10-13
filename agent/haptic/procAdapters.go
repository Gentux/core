package main

import (
	"fmt"

	//nan "nanocloud.com/zeroinstall/lib/libnan"
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

	fmt.Println("TODO GetUsers")

	// var UsersJsonArray []string

	// e := g_pDb.View(func(tx *bolt.Tx) error {

	// 	bucket := tx.Bucket([]byte("users"))
	// 	if bucket == nil {
	// 		return errors.New("Bucket 'users' doesn't exist")
	// 	}

	// 	c := bucket.Cursor()

	// 	for k, v := c.First(); k != nil; k, v = c.Next() {
	// 		UsersJsonArray = append(UsersJsonArray, string(v))
	// 	}

	// 	return nil
	// })

	// return UsersJsonArray, e

	return nil, nil
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
