package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthGmailConfig *oauth2.Config
	oAuthState       = "test"
)

//UserContent models the user content data
type UserContent struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}

func init() {
	oauthGmailConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8181/userprofile",
		ClientID:     "259886612752-hikn6rmnashmm18u4r4js4g00t0hkceb.apps.googleusercontent.com",
		ClientSecret: "CtiASpDed0u8d4rd9lfX7xxE",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/userprofile", handleGoogleCallback)
	fmt.Println(http.ListenAndServe(":8181", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthGmailConfig.AuthCodeURL(oAuthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	callBackPage(w, r, content)
}

func getUserInfo(state string, code string) (*UserContent, error) {
	if state != oAuthState {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := oauthGmailConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	var resp UserContent

	err = json.Unmarshal(contents, &resp)

	fmt.Println("ID :" + resp.ID)
	fmt.Println("Email :" + resp.Email)
	fmt.Println("Link :" + resp.Link)
	fmt.Println("Picture :" + resp.Picture)
	fmt.Println("VerifiedEmail :" + strconv.FormatBool(resp.VerifiedEmail))
	fmt.Println("HD :" + resp.HD)

	//string(resp.VerifiedEmail))

	return &resp, nil
}

func callBackPage(w http.ResponseWriter, r *http.Request, content *UserContent) {

	tmpl, err := template.ParseFiles("userprofile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
