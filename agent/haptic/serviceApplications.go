package main

import (
	"net/http"
)

var ()

type ServiceApplications struct {
}

type GetApplicationsListReply struct {
	ApplicationsJsonArray string
}

func (p *ServiceApplications) GetList(r *http.Request, args *NoArgs, reply *GetApplicationsListReply) error {

	applicationsJson, _ := adapter.GetApplications()
	reply.ApplicationsJsonArray = applicationsJson

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


	adapter.UnpublishApp(args.ApplicationName)
	reply.Result = true

	return nil
}
