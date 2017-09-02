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

// Package session implements a session store for the Go-OAuth examples.  A
// real application should not use this package.
package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

var (
	mu       sync.Mutex
	sessions = make(map[string]map[string]interface{})
)

// Get returns the session data for the request client.
func Get(r *http.Request) (s map[string]interface{}) {
	if c, _ := r.Cookie("session"); c != nil && c.Value != "" {
		mu.Lock()
		s = sessions[c.Value]
		mu.Unlock()
	}
	if s == nil {
		s = make(map[string]interface{})
	}
	return s
}

// Save saves session for the request client.
func Save(w http.ResponseWriter, r *http.Request, s map[string]interface{}) error {
	key := ""
	if c, _ := r.Cookie("session"); c != nil {
		key = c.Value
	}
	if len(s) == 0 {
		if key != "" {
			mu.Lock()
			delete(sessions, key)
			mu.Unlock()
		}
		return nil
	}
	if key == "" {
		var buf [16]byte
		_, err := rand.Read(buf[:])
		if err != nil {
			return err
		}
		key = hex.EncodeToString(buf[:])
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Path:     "/",
			HttpOnly: true,
			Value:    key,
		})
	}
	mu.Lock()
	sessions[key] = s
	mu.Unlock()
	return nil
}
