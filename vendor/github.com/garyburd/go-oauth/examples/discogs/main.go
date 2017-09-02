// Copyright 2014 Gary Burd
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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/garyburd/go-oauth/examples/session"
	"github.com/garyburd/go-oauth/oauth"
)

// Session state keys.
const (
	tempCredKey  = "tempCred"
	tokenCredKey = "tokenCred"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.discogs.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://www.discogs.com/oauth/authorize",
	TokenRequestURI:               "https://api.discogs.com/oauth/access_token",
	Header:                        http.Header{"User-Agent": {"ExampleDiscogsClient/1.0"}},
}

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

func readCredentials() error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &oauthClient.Credentials)
}

// serveLogin gets the OAuth temp credentials and redirects the user to the
// Discogs' authorization page.
func serveLogin(w http.ResponseWriter, r *http.Request) {
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

// serveLogout clears the token credentials
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

// getJSON gets a resource from the Discogs API server and decodes the result as JSON.
func getJSON(cred *oauth.Credentials, endpoint string, form url.Values, v interface{}) error {
	resp, err := oauthClient.Get(nil, cred, "https://api.discogs.com"+endpoint, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
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
		data := struct{ Host string }{r.Host}
		respond(w, homeLoggedOutTmpl, &data)
	} else {
		var data struct {
			Username string
		}
		if err := getJSON(cred, "/oauth/identity", nil, &data); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		respond(w, homeTmpl, &data)
	}
}

var httpAddr = flag.String("addr", ":8080", "HTTP server address")

func main() {
	flag.Parse()
	if err := readCredentials(); err != nil {
		log.Fatalf("Error reading configuration, %v", err)
	}

	http.Handle("/", &authHandler{handler: serveHome, optional: true})
	http.HandleFunc("/login", serveLogin)
	http.HandleFunc("/logout", serveLogout)
	http.HandleFunc("/callback", serveOAuthCallback)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}

var (
	homeLoggedOutTmpl = template.Must(template.New("loggedout").Parse(
		`<html>
<body>
<a href="/login">login</a>
</body>
</html>`))

	homeTmpl = template.Must(template.New("home").Parse(
		`<html>
<body>
<p>Welcome {{.Username}}!
</body></html>`))
)
