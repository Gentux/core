package main

import (
	"log"
	"net/http"
)

var ()

type ServiceUsers struct {
}

// ====================================================================================================

type GetUsersListReply struct {
	Users []User
}

func (p *ServiceUsers) GetList(r *http.Request, args *NoArgs, reply *GetUsersListReply) error {

	if users, err := g_Db.GetUsers(); err != nil {
		LogErrorCode(err)
	} else {

		reply.Users = users

		log.Println(users)
	}

	return nil
}

// ====================================================================================================

func (p *ServiceUsers) GetUser(r *http.Request, args *NoArgs, reply *DefaultReply) error {
	log.Println("TODO GetUser (not sure it need to be done though)")
	return nil
}

// ====================================================================================================

type RegisterUserParam struct {
	Firstname, Lastname, Email, Password string
}

func (p *ServiceUsers) RegisterUser(r *http.Request, args *RegisterUserParam, reply *DefaultReply) error {

	err := adapter.RegisterUser(args.Firstname, args.Lastname, args.Email, args.Password)
	if err != nil {
		reply.Result = false
	} else {
		reply.Result = true
	}

	return nil
}

func (p *ServiceUsers) UpdateUser(r *http.Request, args *RegisterUserParam, reply *DefaultReply) error {

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

	err := adapter.DeleteUser(args.Email)
	if err != nil {
		reply.Result = false
	} else {
		reply.Result = true
	}

	return nil
}
