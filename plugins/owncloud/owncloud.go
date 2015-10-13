package main

import (
	"fmt"
	"io/ioutil"

	"log"

	"net/http"
	"net/url"
)

var (
	protocol, hostname, adminLogin, password string

	g_apiUrl string
)

func Configure(_protocol, _hostname, _adminLogin, _password string) {
	protocol = _protocol
	hostname = _hostname
	adminLogin = _adminLogin
	password = _password

	g_apiUrl = fmt.Sprintf("%s://%s:%s@%s/ocs/v1.php/cloud/users", protocol, adminLogin, password, hostname)
}

func CreateUser(_username, _password string) error {

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

func DeleteUser(_userId string) error {

	deleteUrl := g_apiUrl + "/" + _userId

	if req, err := http.NewRequest("DELETE", deleteUrl, nil); err != nil {
		log.Printf("error creating delete request: %s", err)
		return err
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		log.Printf("error sending request to owncloud: %s", err)
		return err
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Printf("error parsing response from owncloud: %s", err)
		return err
	} else {
		resp.Body.Close()
		log.Printf("owncloud request result body: %s", body)
	}

	return nil
}

func UpdateUserEmail() {

}

func UpdateUserPassword() {

}

// func ListUsers() {

// }

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
