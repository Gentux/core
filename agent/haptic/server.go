package main

import (
	"log"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"net/http"
)

type NoArgs struct {
}

type DefaultReply struct {
	Result  bool
	Code    int
	Message string
}

func StaticHandler(w http.ResponseWriter, pRequest *http.Request) {

	url := pRequest.URL.String()

	if url == "/" {
		url = "/index.html"
	}

	LocalPath := nan.Config().Proxy.FrontendRootDir + url

	http.ServeFile(w, pRequest, LocalPath)
}

func RunServer() {

	// Setup basic HTTP server to serve static content
	http.HandleFunc("/", StaticHandler)

	// Setup RPC server
	pRpcServer := rpc.NewServer()
	pRpcServer.RegisterCodec(json.NewCodec(), "application/json")
	pRpcServer.RegisterService(new(ServiceUsers), "")
	pRpcServer.RegisterService(new(ServiceApplications), "")

	http.Handle("/rpc", pRpcServer)

	log.Println("Now listening on http://localhost:8081")
	e := http.ListenAndServe(":8081", nil)
	if e != nil {
		log.Fatal(e)
	}
}
