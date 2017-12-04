// Copyright 2011 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/examples/session"
	"github.com/garyburd/go-oauth/oauth"
)

// Session state keys.
const (
	tempCredKey  = "tempCred"
	tokenCredKey = "tokenCred"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var signinOAuthClient oauth.Client

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &oauthClient.Credentials)
}

// serveSignin gets the OAuth temp credentials and redirects the user to the
// Twitter's authentication page.
func serveSignin(w http.ResponseWriter, r *http.Request) {
	callback := "http://" + r.Host + "/callback"
	tempCred, err := signinOAuthClient.RequestTemporaryCredentials(nil, callback, nil)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return
	}
	s := session.Get(r)
	s[tempCredKey] = tempCred
	if err := session.Save(w, r, s); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, signinOAuthClient.AuthorizationURL(tempCred, nil), 302)
}

// serveAuthorize gets the OAuth temp credentials and redirects the user to the
// Twitter's authorization page.
func serveAuthorize(w http.ResponseWriter, r *http.Request) {
	callback := "http://" + r.Host + "/callback"
	tempCred, err := oauthClient.RequestTemporaryCredentials(nil, callback, nil)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return
	}
	s := session.Get(r)
	s[tempCredKey] = tempCred
	if err := session.Save(w, r, s); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, oauthClient.AuthorizationURL(tempCred, nil), 302)
}

// serveOAuthCallback handles callbacks from the OAuth server.
func serveOAuthCallback(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	tempCred, _ := s[tempCredKey].(*oauth.Credentials)
	if tempCred == nil || tempCred.Token != r.FormValue("oauth_token") {
		http.Error(w, "Unknown oauth_token.", 500)
		return
	}
	tokenCred, _, err := oauthClient.RequestToken(nil, tempCred, r.FormValue("oauth_verifier"))
	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}
	delete(s, tempCredKey)
	s[tokenCredKey] = tokenCred
	if err := session.Save(w, r, s); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", 302)
}

// serveLogout clears the authentication cookie.
func serveLogout(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	delete(s, tokenCredKey)
	if err := session.Save(w, r, s); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", 302)
}

// authHandler reads the auth cookie and invokes a handler with the result.
type authHandler struct {
	handler  func(w http.ResponseWriter, r *http.Request, c *oauth.Credentials)
	optional bool
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cred, _ := session.Get(r)[tokenCredKey].(*oauth.Credentials)
	if cred == nil && !h.optional {
		http.Error(w, "Not logged in.", 403)
		return
	}
	h.handler(w, r, cred)
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func apiGet(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Get(nil, cred, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func apiPost(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Post(nil, cred, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		p, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("get %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

// respond responds to a request by executing the html template t with data.
func respond(w http.ResponseWriter, t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Print(err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request, cred *oauth.Credentials) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if cred == nil {
		respond(w, homeLoggedOutTmpl, nil)
	} else {
		respond(w, homeTmpl, nil)
	}
}

func serveTimeline(w http.ResponseWriter, r *http.Request, cred *oauth.Credentials) {
	var timeline []map[string]interface{}
	if err := apiGet(
		cred,
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		url.Values{"include_entities": {"true"}},
		&timeline); err != nil {
		http.Error(w, "Error getting timeline, "+err.Error(), 500)
		return
	}
	respond(w, timelineTmpl, timeline)
}

func serveMessages(w http.ResponseWriter, r *http.Request, cred *oauth.Credentials) {
	var dms []map[string]interface{}
	if err := apiGet(
		cred,
		"https://api.twitter.com/1.1/direct_messages.json",
		nil,
		&dms); err != nil {
		http.Error(w, "Error getting timeline, "+err.Error(), 500)
		return
	}
	respond(w, messagesTmpl, dms)
}

func serveFollow(w http.ResponseWriter, r *http.Request, cred *oauth.Credentials) {
	var profile map[string]interface{}
	if err := apiPost(
		cred,
		"https://api.twitter.com/1.1/friendships/create.json",
		url.Values{"screen_name": {"gburd"}, "follow": {"true"}},
		&profile); err != nil {
		http.Error(w, "Error following, "+err.Error(), 500)
		return
	}
	respond(w, followTmpl, profile)
}

var httpAddr = flag.String("addr", ":8080", "HTTP server address")

func main() {
	flag.Parse()
	if err := readCredentials(); err != nil {
		log.Fatalf("Error reading configuration, %v", err)
	}

	// Use a different auth URL for "Sign in with Twitter."
	signinOAuthClient = oauthClient
	signinOAuthClient.ResourceOwnerAuthorizationURI = "https://api.twitter.com/oauth/authenticate"

	http.Handle("/", &authHandler{handler: serveHome, optional: true})
	http.Handle("/timeline", &authHandler{handler: serveTimeline})
	http.Handle("/messages", &authHandler{handler: serveMessages})
	http.Handle("/follow", &authHandler{handler: serveFollow})
	http.HandleFunc("/signin", serveSignin)
	http.HandleFunc("/authorize", serveAuthorize)
	http.HandleFunc("/logout", serveLogout)
	http.HandleFunc("/callback", serveOAuthCallback)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}

var (
	homeLoggedOutTmpl = template.Must(template.New("loggedout").Parse(
		`<html>
<head>
</head>
<body>
<a href="/authorize">Authorize</a> or
<a href="/signin"><img src="http://g.twimg.com/dev/sites/default/files/images_documentation/sign-in-with-twitter-gray.png"></a>
</body>
</html>`))

	homeTmpl = template.Must(template.New("home").Parse(
		`<html>
<head>
</head>
<body>
<p><a href="/timeline">timeline</a>
<p><a href="/messages">direct messages</a>
<p><a href="/follow">follow @gburd</a>
<p><a href="/logout">logout</a>
</body></html>`))

	messagesTmpl = template.Must(template.New("messages").Parse(
		`<html>
<head>
</head>
<body>
<p><a href="/">home</a>
{{range .}}
<p><b>{{.sender.name}}</b> {{.text}}
{{end}}
</body></html>`))

	timelineTmpl = template.Must(template.New("timeline").Parse(
		`<html>
<head>
</head>
<body>
<p><a href="/">home</a>
{{range .}}
<p><b>{{.user.name}}</b> {{.text}}
{{with .entities}}
    {{with .urls}}<br><i>urls:</i> {{range .}}{{.expanded_url}}{{end}}{{end}}
    {{with .hashtags}}<br><i>hashtags:</i> {{range .}}{{.text}}{{end}}{{end}}
    {{with .user_mentions}}<br><i>user_mentions:</i> {{range .}}{{.screen_name}}{{end}}{{end}}
{{end}}
{{end}}
</body></html>`))

	followTmpl = template.Must(template.New("follow").Parse(
		`<html>
<head>
</head>
<body>
<p><a href="/">home</a>
<p>You are now following <a href="https://twitter.com/{{.screen_name}}">{{.name}}</a>
</body></html>`))
)
