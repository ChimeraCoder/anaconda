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

package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/garyburd/go-oauth/oauth"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

// context stores context associated with an HTTP request.
type Context struct {
	c context.Context
	r *http.Request
	w http.ResponseWriter
	u *user.User
}

// handler adapts a function to an http.Handler
type handler func(c *Context) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	c := Context{
		c: ctx,
		r: r,
		w: w,
		u: user.Current(ctx),
	}

	if c.u == nil {
		url, _ := user.LoginURL(c.c, c.r.URL.Path)
		http.Redirect(w, r, url, 301)
		return
	}

	err := h(&c)
	if err != nil {
		http.Error(w, "server error", 500)
		log.Errorf(c.c, "error %v", err)
	}
}

// userInfo is stored in the App Engine datastore with key email.
type userInfo struct {
	TwitterCred oauth.Credentials
}

// getUserInfo returns information about the currently logged in user.
func (c *Context) getUserInfo() (*userInfo, error) {
	key := datastore.NewKey(c.c, "user", c.u.Email, 0, nil)
	var u userInfo
	err := datastore.Get(c.c, key, &u)
	if err == datastore.ErrNoSuchEntity {
		err = nil
	}
	return &u, err
}

// updateUserInfo updates information about the currently logged in user.
func (c *Context) updateUserInfo(f func(u *userInfo)) error {
	key := datastore.NewKey(c.c, "user", c.u.Email, 0, nil)
	return datastore.RunInTransaction(c.c, func(ctx context.Context) error {
		var u userInfo
		err := datastore.Get(ctx, key, &u)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		f(&u)
		_, err = datastore.Put(ctx, key, &u)
		return err
	}, nil)
}

type connectInfo struct {
	Secret   string
	Redirect string
}

// serveTwitterConnect gets the OAuth temp credentials and redirects the user to the
// Twitter's authorization page.
func serveTwitterConnect(c *Context) error {
	httpClient := urlfetch.Client(c.c)
	callback := "http://" + c.r.Host + "/twitter/callback"
	tempCred, err := oauthClient.RequestTemporaryCredentials(httpClient, callback, nil)
	if err != nil {
		return err
	}

	ci := connectInfo{Secret: tempCred.Secret, Redirect: c.r.FormValue("redirect")}
	err = memcache.Gob.Set(c.c, &memcache.Item{Key: tempCred.Token, Object: &ci})
	if err != nil {
		return err
	}
	http.Redirect(c.w, c.r, oauthClient.AuthorizationURL(tempCred, nil), 302)
	return nil
}

// serveTwitterCallback handles callbacks from the Twitter OAuth server.
func serveTwitterCallback(c *Context) error {
	token := c.r.FormValue("oauth_token")
	var ci connectInfo
	_, err := memcache.Gob.Get(c.c, token, &ci)
	if err != nil {
		return err
	}
	memcache.Delete(c.c, token)
	tempCred := &oauth.Credentials{
		Token:  token,
		Secret: ci.Secret,
	}

	httpClient := urlfetch.Client(c.c)
	tokenCred, _, err := oauthClient.RequestToken(httpClient, tempCred, c.r.FormValue("oauth_verifier"))
	if err != nil {
		return err
	}

	if err := c.updateUserInfo(func(u *userInfo) { u.TwitterCred = *tokenCred }); err != nil {
		return err
	}
	http.Redirect(c.w, c.r, ci.Redirect, 302)
	return nil
}

// serveTwitterDisconnect clears the user's Twitter credentials.
func serveTwitterDisconnect(c *Context) error {
	if err := c.updateUserInfo(func(u *userInfo) { u.TwitterCred = oauth.Credentials{} }); err != nil {
		return err
	}
	http.Redirect(c.w, c.r, c.r.FormValue("redirect"), 302)
	return nil
}

func serveHome(c *Context) error {
	if c.r.URL.Path != "/" {
		http.NotFound(c.w, c.r)
		return nil
	}
	u, err := c.getUserInfo()
	if err != nil {
		return err
	}

	var data = struct {
		Connected bool
		Timeline  []map[string]interface{}
	}{
		Connected: u.TwitterCred.Token != "" && u.TwitterCred.Secret != "",
	}

	if data.Connected {
		httpClient := urlfetch.Client(c.c)
		resp, err := oauthClient.Get(httpClient, &u.TwitterCred, "https://api.twitter.com/1.1/statuses/home_timeline.json", nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			p, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("get %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
		}
		if err := json.NewDecoder(resp.Body).Decode(&data.Timeline); err != nil {
			return err
		}
	}

	c.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return homeTmpl.Execute(c.w, &data)
}

func init() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &oauthClient.Credentials); err != nil {
		panic(err)
	}
	http.Handle("/", handler(serveHome))
	http.Handle("/twitter/connect", handler(serveTwitterConnect))
	http.Handle("/twitter/disconnect", handler(serveTwitterDisconnect))
	http.Handle("/twitter/callback", handler(serveTwitterCallback))
}

var homeTmpl = template.Must(template.New("home").Parse(
	`<html>
<head>
</head>
<body>
{{if .Connected}}
    <a href="/twitter/disconnect?redirect=/">Disconnect Twitter account</a>
    {{range .Timeline}}
    <p><b>{{html .user.name}}</b> {{html .text}}
    {{end}}
{{else}}
    <a href="/twitter/connect?redirect=/">Connect Twitter account</a>
{{end}}
</body>
</html>`))
