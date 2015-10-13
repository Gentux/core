package main

import (
	"fmt"
	"io/ioutil"

	"log"

	"net/http"
	"net/url"
)

var (
	Protocol, Url, Login, Password string

	g_apiUrl string
)

func OwncloudConfigure(configParams ConfigureParams) error {

	Protocol = configParams.Protocol
	Url = configParams.Url
	Login = configParams.Login
	Password = configParams.Password

	g_apiUrl = fmt.Sprintf("%s://%s:%s@%s/ocs/v1.php/cloud/users", Protocol, Login, Password, Url)

	return nil
}

func OwncloudAddUser(_username, _password string) error {

	values := make(url.Values)
	values.Set("userid", _username)
	values.Set("password", _password)

	r, err := http.PostForm(g_apiUrl, values)
	if err != nil {
		log.Printf("error posting to owncloud: %s", err)
		return err
	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Printf("error reading response body: %s", err)
		return err
	}
	r.Body.Close()
	log.Printf("owncloud post result body: %s", body)

	return nil
}

func OwncloudDeleteUser(_username string) error {

	url := g_apiUrl + "/" + _username

	req, err := http.NewRequest("DELETE", url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Owncloud plugin error %s when invoking url: %s\n", err.Error(), url)
	}

	if resp != nil {
		resp.Body.Close()
	}

	return nil
}

func UpdateUserEmail() {

}

func UpdateUserPassword() {

}

// func main() {

// "application/x-www-form-urlencoded"

// GET http://admin:secret@example.com/ocs/v1.php/cloud/users/Frank
// Returns information on the user Frank

// GET http://admin:secret@example.com/ocs/v1.php/cloud/users?search=Frank
// Returns list of users matching the search string.

// PUT http://admin:secret@example.com/ocs/v1.php/cloud/users/Frank -d key="email" -d value="franksnewemail@example.org"
// Updates the email address for the user Frank

// Works also for "password". However, it is not recommended to use the API in this case, since ownCloud is encrypted (the decryption key should be used).

//}
