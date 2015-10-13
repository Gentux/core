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
