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
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

func set(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// get the parameters k and v from the request
	key := r.FormValue("k")
	value := r.FormValue("v")

	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: 1 * time.Hour,
	}

	err := memcache.Set(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := r.FormValue("k")

	item, err := memcache.Get(ctx, key)
	switch err {
	case nil:
		fmt.Fprintf(w, "%s", item.Value)
	case memcache.ErrCacheMiss:
		fmt.Fprint(w, "key not found")
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
}
