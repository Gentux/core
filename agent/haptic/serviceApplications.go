package main

import (
	"errors"
	"net/http"
)

var ()

type ServiceApplications struct {
}

type GetApplicationsListReply struct {
	Applications []Connection
}

func (p *ServiceApplications) GetList(r *http.Request, args *NoArgs, reply *GetApplicationsListReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	connections, _ := adapter.GetApplications()
	reply.Applications = connections

	return nil
}

type SamAccountParam struct {
	Sam string
}

func (p *ServiceApplications) GetListForSamAccount(r *http.Request, args *SamAccountParam, reply *GetApplicationsListReply) error {

	connections, _ := adapter.GetApplicationsForSamAccount(args.Sam)
	reply.Applications = connections

	return nil
}

// ====================================================================================================

func (p *ServiceApplications) GetApplication(r *http.Request, args *NoArgs, reply *DefaultReply) error {
	return nil
}

// ====================================================================================================

type RegisterApplicationParam struct {
	ApplicationName, ApplicationPath string
}

func (p *ServiceApplications) RegisterApplication(r *http.Request, args *RegisterApplicationParam, reply *DefaultReply) error {

	// TODO Implement this

	return nil
}

// ====================================================================================================

func (p *ServiceApplications) UpdateApplication(r *http.Request, args *RegisterApplicationParam, reply *DefaultReply) error {

	// TODO Implement this

	return nil
}

// ====================================================================================================

type UnpublishApplicationParam struct {
	ApplicationName string
}

func (p *ServiceApplications) UnpublishApplication(r *http.Request, args *UnpublishApplicationParam, reply *DefaultReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	adapter.UnpublishApp(args.ApplicationName)
	reply.Result = true

	return nil
}
