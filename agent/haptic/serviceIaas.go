package main

import (
	"net/http"
)

var ()

type ServiceIaas struct {
}

type GetVmListReply struct {
	VmListJsonArray string
}

func (p *ServiceIaas) GetList(r *http.Request, args *NoArgs, reply *GetVmListReply) error {

	vmListJson, _ := adapter.GetVmList()
	reply.VmListJsonArray = vmListJson

	return nil
}

type RequestState struct {
	Success bool
}

func (p *ServiceIaas) Download(r *http.Request, args *NoArgs, reply *RequestState) error {

	requestState, _ := adapter.DownloadWindowsVm()
	reply.Success = requestState

	return nil
}

type VmNameArgs struct {
	VmName string
}

func (p *ServiceIaas) Start(r *http.Request, args *VmNameArgs, reply *RequestState) error {

	requestState, _ := adapter.StartVm(args.VmName)
	reply.Success = requestState

	return nil
}
