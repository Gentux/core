package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var ()

type ServiceUsers struct {
}

// ====================================================================================================

type GetUsersListReply struct {
	UsersJsonArray string
}

func (p *ServiceUsers) GetList(r *http.Request, args *NoArgs, reply *GetUsersListReply) error {

	var (
		users      []map[string]string
		users_mail []string
		e          error
	)

	users_mail, e = adapter.GetUsers()
	for _, user_mail := range users_mail {
		user := map[string]string{}

		user["Firstname"], user["Lastname"], user["Email"], user["Password"], user["License"] =
			GetUserAccountParamsForActivation(user_mail)

		users = append(users, user)
	}

	if e != nil {
		log.Println(e)
	}

	message, _ := json.Marshal(users)
	reply.UsersJsonArray = string(message)
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

	adapter.DeleteUser(args.Email)

	reply.Result = true

	return nil
}
