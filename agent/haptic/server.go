package main

import (
	"log"
	"net/http"
	"time"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

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

func SecureHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("nanocloud")
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		value := make(map[string]string)
		err = cookieHandler.Decode("nanocloud", cookie.Value, &value)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		//		isUserRegistered, _ := g_Db.IsUserRegistered(value["email"])
		//		expirationTime, _ := time.Parse(time.RFC3339, value["expirationTime"])
		//		if !isUserRegistered || time.Now() >= expirationTime {
		// TODO Return 401
		//		}

		h.ServeHTTP(w, r)
	})
}

func loginHandler(response http.ResponseWriter, request *http.Request) {

	// TODO Reject all request not in "POST"
	var (
		email          string = request.FormValue("email")
		password       string = request.FormValue("password")
		user           User
		redirectTarget string = "/"
		expirationTime time.Time
	)

	user, _ = g_Db.GetUser(email)
	if password != user.Password {
		redirectTarget = "/login.html"
	} else {
		expirationTime = time.Now().Add(4 * time.Hour)
		value := map[string]string{
			"email":          user.Email,
			"expirationTime": expirationTime.Format(time.RFC3339),
		}
		user.TokenExpirationTime = user.CreationTime
		if encoded, err := cookieHandler.Encode("nanocloud", value); err == nil {
			cookie := &http.Cookie{
				Name:  "nanocloud",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(response, cookie)
		}
		redirectTarget = "/"
	}

	http.Redirect(response, request, redirectTarget, 302)
}

func RunServer() {

	// Setup basic HTTP server to serve static content
	http.HandleFunc("/", StaticHandler)

	// Login handler
	http.HandleFunc("/login", loginHandler)

	// Setup RPC server
	pRpcServer := rpc.NewServer()
	pRpcServer.RegisterCodec(json.NewCodec(), "application/json")
	pRpcServer.RegisterService(new(ServiceAuthentication), "")
	pRpcServer.RegisterService(new(ServiceIaas), "")
	pRpcServer.RegisterService(new(ServiceApplications), "")
	pRpcServer.RegisterService(new(ServiceUsers), "")

	secureHandler := SecureHandler(pRpcServer)
	http.Handle("/rpc", secureHandler)

	Log("Now listening on http://localhost:8081")
	e := http.ListenAndServe(":8081", nil)
	if e != nil {
		log.Fatal(e)
	}
}
