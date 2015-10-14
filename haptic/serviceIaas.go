package main

import (
	"errors"
	"net/http"
)

var ()

type ServiceIaas struct {
}

type GetVmListReply struct {
	VmListJsonArray string
}

func (p *ServiceIaas) GetList(r *http.Request, args *NoArgs, reply *GetVmListReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	vmListJson, _ := adapter.GetVmList()
	reply.VmListJsonArray = vmListJson

	return nil
}

type RequestState struct {
	Success bool
}

func (p *ServiceIaas) Download(r *http.Request, args *NoArgs, reply *RequestState) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	requestState, _ := adapter.DownloadWindowsVm()
	reply.Success = requestState

	return nil
}

func (p *ServiceIaas) DownloadStatus(r *http.Request, args *NoArgs, reply *RequestState) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	requestState, _ := adapter.DownloadStatus()
	reply.Success = requestState

	return nil
}

type VmNameArgs struct {
	VmName string
}

func (p *ServiceIaas) Start(r *http.Request, args *VmNameArgs, reply *RequestState) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	requestState, _ := adapter.StartVm(args.VmName)
	reply.Success = requestState

	return nil
}
