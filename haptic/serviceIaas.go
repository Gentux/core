/*
 * Nanocloud community -- transform any application into SaaS solution
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
