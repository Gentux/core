package main

import (
	"encoding/json"
	"log"
	"net/http"

	nan "nanocloud.com/zeroinstall/lib/libnan"
)

var ()

type ServiceUsers struct {
}

// ====================================================================================================

type UserInfo struct {
	Firstname, Lastname, Email, Password string
}

// ====================================================================================================

type GetUsersListReply struct {
	UsersJsonArray string
}

//DESIGN note : if users list tends to become huge we'll probably have to break it down in subsets
func (p *ServiceUsers) GetList(r *http.Request, args *NoArgs, reply *GetUsersListReply) error {

	usersEmails, e := ListUserEmails()

	if e != nil {
		LogError("ServiceUsers.GetList got error from ListUserEmails: %s", e.Message)
		return nil
	} else if len(usersEmails) == 0 {
		return nil
	}

	userInfos := make([]UserInfo, len(usersEmails), len(usersEmails))

	log.Printf("ServiceUsers.GetList : returning list of %d users\n", len(usersEmails))

	for idx, userEmail := range usersEmails {

		var e *nan.Error

		userInfos[idx].Firstname, userInfos[idx].Lastname, userInfos[idx].Email, userInfos[idx].Password,
			_, e =
			GetUserAccountParamsForActivation(userEmail)

		if e != nil {
			nan.LogErrorCode(e)
		}
	}

	if message, e := json.Marshal(userInfos); e != nil {
		LogError("ServiceUsers.GetList failed to marshall array of %d userInfos with error: %s",
			len(userInfos), e.Error())
	} else {
		reply.UsersJsonArray = string(message)
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

	adapter.RegisterUser(args.Firstname, args.Lastname, args.Email, args.Password)

	return nil
}

func (p *ServiceUsers) ActivateUser(r *http.Request, args *RegisterUserParam, reply *DefaultReply) error {

	err := adapter.ActivateUser(args.Email)

	reply.Code = err.Code
	reply.Message = err.Message

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

	adapter.DeleteUser(args.Email)

	reply.Result = true

	return nil
}
