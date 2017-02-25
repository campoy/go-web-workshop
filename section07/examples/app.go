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

func getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := datastore.NewKey(ctx, "Person", "gopher", 0, nil)

	var p Person
	err := datastore.Get(ctx, key, &p)
	if err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, p)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var p []Person

	// create a new query on the kind Person
	q := datastore.NewQuery("Person")

	// select only values where field Age is 10 or lower
	q = q.Filter("Age <=", 10)

	// order all the values by the Name field
	q = q.Order("Name")

	// and finally execute the query retrieving all values into p.
	_, err := q.GetAll(ctx, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, p)
}

func chainedQueryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var p []Person

	// create a new query on the kind Person
	q := datastore.NewQuery("Person").
		Filter("Age <=", 10).
		Order("Name")

	// and finally execute the query retrieving all values into p.
	_, err := q.GetAll(ctx, &p)
	if err != nil {
		// handle the error
	}
}

func init() {
	http.HandleFunc("/complete", completeHandler)
	http.HandleFunc("/incomplete", incompleteHandler)
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/chainedQuery", chainedQueryHandler)
	http.HandleFunc("/get", getHandler)
}
