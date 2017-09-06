// Copyright 2015 Gary Burd
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
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

type client struct {
	client oauth.Client
	token  oauth.Credentials
}

func (c *client) get(urlStr string, params url.Values, v interface{}) error {
	resp, err := c.client.Get(nil, &c.token, urlStr, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("yelp status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

func readCredentials(c *client) error {
	b, err := ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}
	var creds struct {
		ConsumerKey    string
		ConsumerSecret string
		Token          string
		TokenSecret    string
	}
	if err := json.Unmarshal(b, &creds); err != nil {
		return err
	}
	c.client.Credentials.Token = creds.ConsumerKey
	c.client.Credentials.Secret = creds.ConsumerSecret
	c.token.Token = creds.Token
	c.token.Secret = creds.TokenSecret
	return nil
}

func main() {
	var c client
	if err := readCredentials(&c); err != nil {
		log.Fatal(err)
	}

	var data struct {
		Businesses []struct {
			Name     string
			Location struct {
				DisplayAddress []string `json:"display_address"`
			}
		}
	}
	form := url.Values{"term": {"food"}, "location": {"San Francisco"}}
	if err := c.get("http://api.yelp.com/v2/search", form, &data); err != nil {
		log.Fatal(err)
	}

	for _, b := range data.Businesses {
		addr := ""
		if len(b.Location.DisplayAddress) > 0 {
			addr = b.Location.DisplayAddress[0]
		}
		log.Println(b.Name, addr)
	}
}
