/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
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
	"fmt"
	"net/http"
)

type ServiceHistory struct {
}

type GetHistoryListReply struct {
	Histories []History
}

type HistoryArgs struct {
	ConnectionId string
	StartDate    string
	EndDate      string
}

func (p *ServiceHistory) GetList(r *http.Request, args *NoArgs, reply *GetHistoryListReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	reply.Histories = GetHistory()

	return nil
}

func (p *ServiceHistory) GetListForUser(r *http.Request, args *NoArgs, reply *GetHistoryListReply) error {

	value := make(map[string]string)
	cookie, _ := r.Cookie("nanocloud")
	cookieHandler.Decode("nanocloud", cookie.Value, &value)
	user, _ := g_Db.GetUser(value["email"])

	reply.Histories = GetHistoryForUser(user.Email)
	return nil
}

func (p *ServiceHistory) Add(r *http.Request, args *HistoryArgs, reply *RequestState) error {

	fmt.Println(*args)
	reply.Success = AddHistory(*args)
	return nil
}
