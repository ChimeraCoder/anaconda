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
	companyKey   = "company"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://oauth.intuit.com/oauth/v1/get_request_token",
	ResourceOwnerAuthorizationURI: "https://appcenter.intuit.com/Connect/Begin",
	TokenRequestURI:               "https://oauth.intuit.com/oauth/v1/get_access_token",
	Header:                        http.Header{"Accept": {"application/json"}},
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
// Quickbooks's authorization page.
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
	s[companyKey] = r.FormValue("realmId")
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
	delete(s, companyKey)
	if err := session.Save(w, r, s); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", 302)
}

// authHandler reads the auth cookie and invokes a handler with the result.
type authHandler struct {
	handler  func(w http.ResponseWriter, r *http.Request, c *oauth.Credentials, company string)
	optional bool
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	cred, _ := s[tokenCredKey].(*oauth.Credentials)
	if cred == nil && !h.optional {
		http.Error(w, "Not logged in.", 403)
		return
	}
	company, _ := s[companyKey].(string)
	h.handler(w, r, cred, company)
}

func callAPI(cred *oauth.Credentials, company string, endpoint string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Get(nil,
		cred,
		fmt.Sprintf("https://qb.sbfinance.intuit.com/v3/company/%s/%s", company, endpoint),
		form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

func serveHome(w http.ResponseWriter, r *http.Request, cred *oauth.Credentials, company string) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if cred == nil {
		data := struct{ Host string }{r.Host}
		respond(w, homeLoggedOutTmpl, &data)
	} else {
		respond(w, homeTmpl, nil)
	}
}

func serveAccounts(w http.ResponseWriter, r *http.Request, cred *oauth.Credentials, company string) {
	var accounts map[string]interface{}
	if err := callAPI(
		cred,
		company,
		"query",
		url.Values{"query": {"select * from Account"}},
		&accounts); err != nil {
		http.Error(w, "Error getting accounts, "+err.Error(), 500)
		return
	}
	respond(w, accountsTmpl, accounts)
}

var httpAddr = flag.String("addr", ":8080", "HTTP server address")

func main() {
	flag.Parse()
	if err := readCredentials(); err != nil {
		log.Fatalf("Error reading configuration, %v", err)
	}

	http.Handle("/", &authHandler{handler: serveHome, optional: true})
	http.Handle("/accounts", &authHandler{handler: serveAccounts})
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
<head>
<script src="https://js.appcenter.intuit.com/Content/IA/intuit.ipp.anywhere.js" type="text/javascript"></script>
<script>intuit.ipp.anywhere.setup({ menuProxy: '', grantUrl: 'http://{{.Host}}/login'}); </script>
</head>
<body>
<ipp:connectToIntuit></ipp:connectToIntuit>
</body>
</html>`))

	homeTmpl = template.Must(template.New("home").Parse(
		`<html>
<head>
</head>
<body>
<p><a href="/accounts">accounts</a>
</body></html>`))

	accountsTmpl = template.Must(template.New("accounts").Parse(
		`<html>
<head>
</head>
<body>
<p><a href="/">home</a>
{{range .QueryResponse.Account}}<p>{{.Name}}{{end}}
<{{.}}
</body></html>`))
)
