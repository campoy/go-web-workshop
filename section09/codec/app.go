// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

type Person struct {
	Name     string `json:"name"`
	AgeYears int    `json:"age_years"`
}

func set(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := &memcache.Item{
		Key:        "last_person",
		Object:     p, // we set the Object field instead of Value
		Expiration: 1 * time.Hour,
	}

	// we use the JSON codec
	err := memcache.JSON.Set(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var p Person
	_, err := memcache.JSON.Get(ctx, "last_person", &p)
	if err == nil {
		json.NewEncoder(w).Encode(p)
		return
	}
	if err == memcache.ErrCacheMiss {
		fmt.Fprint(w, "key not found")
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func init() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
}
