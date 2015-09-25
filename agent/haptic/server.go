package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func Enforce(profile string, encodedCookie string) bool {
	// TODO method calling this function should return 403 status code
	value := make(map[string]string)
	cookieHandler.Decode("nanocloud", encodedCookie, &value)

	user, _ := g_Db.GetUser(value["email"])

	return user.Profile == "admin" || profile == user.Profile
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	defer file.Close()

	// Compute a hash name for this file on disk
	tempDst := filepath.Join(nan.Config().CommonBaseDir, "uploads", header.Filename)
	if err = os.MkdirAll(filepath.Join(nan.Config().CommonBaseDir, "uploads"), os.ModePerm); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tmpOutput, err := os.Create(tempDst)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer tmpOutput.Close()

	// write the content from POST to the file
	_, err = io.Copy(tmpOutput, file)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)

	go SyncUploadedFile(tempDst)
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

		isUserRegistered, _ := g_Db.IsUserRegistered(value["email"])
		expirationTime, _ := time.Parse(time.RFC3339, value["expirationTime"])
		if isUserRegistered == false && time.Now().After(expirationTime) {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func loginHandler(response http.ResponseWriter, request *http.Request) {

	// TODO Reject all request not in "POST"
	var (
		email          string = request.FormValue("email")
		password       string = request.FormValue("password")
		user           User
		expirationTime time.Time
	)

	user, err := g_Db.GetUser(email)
	if err != nil || user.Email == "" || password != user.Password {
		http.Error(response, http.StatusText(403), 403)
	} else {
		expirationTime = time.Now().Add(4 * time.Hour)
		value := map[string]string{
			"email":          user.Email,
			"expirationTime": expirationTime.Format(time.RFC3339),
		}
		if encoded, err := cookieHandler.Encode("nanocloud", value); err == nil {
			cookie := &http.Cookie{
				Name:     "nanocloud",
				Value:    encoded,
				Path:     "/",
				Expires:  expirationTime,
				HttpOnly: true,
			}
			http.SetCookie(response, cookie)
		}
	}

	if user.Profile == "admin" {
		http.Redirect(response, request, "/admin.html", 302)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	if _, err := request.Cookie("nanocloud"); err != nil {
		http.Redirect(response, request, "/", 302)
		return
	}

	expirationTime := time.Now()
	value := map[string]string{
		"email":          "",
		"expirationTime": expirationTime.Format(time.RFC3339),
		"expired":        "true",
	}
	if encoded, err := cookieHandler.Encode("nanocloud", value); err == nil {
		cookie := &http.Cookie{
			Name:     "nanocloud",
			Value:    encoded,
			Path:     "/",
			Expires:  expirationTime,
			HttpOnly: true,
		}
		http.SetCookie(response, cookie)
	}
	http.Redirect(response, request, "/", 302)
}

func RunServer() {

	// Setup basic HTTP server to serve static content
	http.HandleFunc("/", StaticHandler)

	// Login handler
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Upload file handler
	http.HandleFunc("/upload", uploadHandler)

	// Setup RPC server
	pRpcServer := rpc.NewServer()
	pRpcServer.RegisterCodec(json.NewCodec(), "application/json")
	pRpcServer.RegisterService(new(ServiceIaas), "")
	pRpcServer.RegisterService(new(ServiceApplications), "")
	pRpcServer.RegisterService(new(ServiceUsers), "")
	pRpcServer.RegisterService(new(ServiceStats), "")

	secureHandler := SecureHandler(pRpcServer)
	http.Handle("/rpc", secureHandler)

	Log("Now listening on http://localhost:8081")
	e := http.ListenAndServe(":8081", nil)
	if e != nil {
		log.Fatal(e)
	}
}
