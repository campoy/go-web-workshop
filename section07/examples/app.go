// Copyright 2016 Google Inc. All rights reserved.
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

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Person contains the name and age of a person.
type Person struct {
	Name     string
	AgeYears int
}

func completeHandler(w http.ResponseWriter, r *http.Request) {
	// create a new App Engine context from the HTTP request.
	ctx := appengine.NewContext(r)

	p := &Person{Name: "gopher", AgeYears: 5}

	// create a new complete key of kind Person and value gopher.
	key := datastore.NewKey(ctx, "Person", "gopher", 0, nil)
	// put p in the datastore.
	key, err := datastore.Put(ctx, key, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "gopher stored with key %v", key)
}

func incompleteHandler(w http.ResponseWriter, r *http.Request) {
	// create a new App Engine context from the HTTP request.
	ctx := appengine.NewContext(r)

	p := &Person{Name: "gopher", AgeYears: 5}

	// create a new complete key of kind Person.
	key := datastore.NewIncompleteKey(ctx, "Person", nil)
	// put p in the datastore.
	key, err := datastore.Put(ctx, key, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "gopher stored with key %v", key)
}

func init() {
	http.HandleFunc("/complete", completeHandler)
	http.HandleFunc("/incomplete", incompleteHandler)
}
