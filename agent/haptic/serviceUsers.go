package main

import (
	// "fmt"
	//"errors"
	"log"
	"net/http"
)

var ()

type ServiceUsers struct {
}

type NoArgs struct {
}

type DefaultReply struct {
	Result bool
}

// ====================================================================================================

type GetUsersListReply struct {
	UsersJsonArray []string
}

func (p *ServiceUsers) GetList(r *http.Request, args *NoArgs, reply *GetUsersListReply) error {

	var e error
	reply.UsersJsonArray, e = adapter.GetUsers()

	if e != nil {
		log.Println(e)
	}

	return nil
}

// ====================================================================================================

func (p *ServiceUsers) GetUser(r *http.Request, args *NoArgs, reply *DefaultReply) error {
	log.Println("GetUser")
	return nil
}

// ====================================================================================================

type RegisterUserParam struct {
	Firstname, Lastname, Email, Password string
}

func (p *ServiceUsers) RegisterUser(r *http.Request, args *RegisterUserParam, reply *DefaultReply) error {

	adapter.RegisterUser(args.Firstname, args.Lastname, args.Email, args.Password)

	return nil
}

// ====================================================================================================

type UpdatePasswordParam struct {
	Email, Password string
}

func (p *ServiceUsers) UpdateUserPassword(r *http.Request, args *UpdatePasswordParam, reply *DefaultReply) error {

	adapter.UpdateUserPassword(args.Email, args.Password)

	return nil
}

// ====================================================================================================

type DeleteUserParam struct {
	Email string
}

func (p *ServiceUsers) DeleteUser(r *http.Request, args *DeleteUserParam, reply *DefaultReply) error {

	adapter.DeleteUser(args.Email)

	return nil
}
