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
	"log"
	"net/http"

	// nan "nanocloud.com/core/lib/libnan"
)

var ()

type Stat struct {
	Activated    bool
	Email        string
	Firstname    string
	Lastname     string
	Password     string
	Sam          string
	CreationTime string
	Profile      string
}

type ServiceStats struct {
}

// ====================================================================================================

type GetStatsListReply struct {
	Stats []Stat
}

//DESIGN note : if stats list tends to become huge we'll probably have to break it down in subsets
func (p *ServiceStats) GetList(r *http.Request, args *NoArgs, reply *GetStatsListReply) error {

	cookie, _ := r.Cookie("nanocloud")
	if Enforce("admin", cookie.Value) == false {
		return errors.New("You need admin permission to perform this action")
	}

	if stats, err := g_Db.GetStats(); err != nil {
		LogErrorCode(err)
	} else {

		reply.Stats = stats

		log.Println(stats)
	}

	return nil
}

// ====================================================================================================

func (p *ServiceStats) GetStat(r *http.Request, args *NoArgs, reply *DefaultReply) error {
	log.Println("TODO GetStat (not sure it need to be done though)")
	return nil
}

// ====================================================================================================

// func (p *ServiceStats) UpdateStat(r *http.Request, args *RegisterStatParam, reply *DefaultReply) error {

// 	cookie, _ := r.Cookie("nanocloud")
// 	if Enforce("admin", cookie.Value) == false {
// 		return errors.New("You need admin permission to perform this action")
// 	}

// 	adapter.RegisterStat(args.Firstname, args.Lastname, args.Email, args.Password)

// 	return nil
// }
