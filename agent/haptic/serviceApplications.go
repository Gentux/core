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

type DeleteApplicationParam struct {
	Email string
}

func (p *ServiceApplications) DeleteApplication(r *http.Request, args *DeleteApplicationParam, reply *DefaultReply) error {

	// TODO Implement this

	return nil
}
