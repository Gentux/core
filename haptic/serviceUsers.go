package main

import (
	"errors"
	"log"
	"net/http"

	// nan "nanocloud.com/core/lib/libnan"
)

var ()

type ServiceUsers struct {
}

// ====================================================================================================

type GetUsersListReply struct {
	Users []User
}

//DESIGN note : if users list tends to become huge we'll probably have to break it down in subsets
func (p *ServiceUsers) GetList(r *http.Request, args *NoArgs, reply *GetUsersListReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

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

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	err := adapter.RegisterUser(args.Firstname, args.Lastname, args.Email, args.Password)
	if err != nil {
		reply.Result = false
	} else {
		reply.Result = true
	}

	return nil
}

func (p *ServiceUsers) ActivateUser(r *http.Request, args *RegisterUserParam, reply *DefaultReply) error {

	err := adapter.ActivateUser(args.Email)

	reply.Code = err.Code
	reply.Message = err.Message

	return nil
}

func (p *ServiceUsers) UpdateUser(r *http.Request, args *RegisterUserParam, reply *DefaultReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	adapter.RegisterUser(args.Firstname, args.Lastname, args.Email, args.Password)

	return nil
}

// ====================================================================================================

type UpdatePasswordParam struct {
	Email, Password string
}

func (p *ServiceUsers) UpdateUserPassword(r *http.Request, args *UpdatePasswordParam, reply *DefaultReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	err := adapter.UpdateUserPassword(args.Email, args.Password)
	if err != nil {
		reply.Result = false
		reply.Code = err.Code
		reply.Message = err.Message
	} else {
		reply.Result = true
	}

	return nil
}

// ====================================================================================================

type DeleteUserParam struct {
	Email string
}

func (p *ServiceUsers) DeleteUser(r *http.Request, args *DeleteUserParam, reply *DefaultReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	err := adapter.DeleteUser(args.Email)
	if err != nil {
		reply.Result = false
	} else {
		reply.Result = true
	}

	return nil
}