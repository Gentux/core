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
