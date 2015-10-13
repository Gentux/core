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

func (p *ServiceApplications) GetListForCurrentUser(r *http.Request, args *NoArgs, reply *GetApplicationsListReply) error {

	value := make(map[string]string)
	cookie, _ := r.Cookie("nanocloud")
	cookieHandler.Decode("nanocloud", cookie.Value, &value)
	user, _ := g_Db.GetUser(value["email"])

	connections, e := adapter.GetApplicationsForSamAccount(user.Sam)
	if e != nil {
		return errors.New("Unable to retrieve connection for Email " + user.Email)
	}

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
