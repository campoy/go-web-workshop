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
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
)

const (
	namespaceKind = "namespace"
	valueKind     = "value"
)

type value struct{ Value string }

func init() {
	r := mux.NewRouter()
	r.Handle("/", csrfHandler{http.HandlerFunc(getNamespaces)}).Methods("GET")
	r.Handle("/{namespace}/", csrfHandler{withNamespace{getAll, false}}).Methods("GET")
	r.Handle("/{namespace}/{key}", csrfHandler{withNamespace{getOne, false}}).Methods("GET")
	r.Handle("/{namespace}/{key}", csrfHandler{withNamespace{put, true}}).Methods("PUT")
	r.Handle("/{namespace}/{key}", csrfHandler{withNamespace{delete, false}}).Methods("DELETE")
	http.Handle("/", r)
}

func getNamespaces(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "getNamespaces")

	keys, err := datastore.NewQuery(namespaceKind).KeysOnly().GetAll(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	for _, key := range keys {
		fmt.Fprintln(w, key.StringID())
	}
}

type csrfHandler struct{ h http.Handler }

func (h csrfHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	h.h.ServeHTTP(w, r)
}

type withNamespace struct {
	h func(w http.ResponseWriter, r *http.Request, namespace *datastore.Key)

	createIfMissing bool
}

func (h withNamespace) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	namespace := mux.Vars(r)["namespace"]
	key := datastore.NewKey(ctx, namespaceKind, namespace, 0, nil)
	if err := datastore.Get(ctx, key, new(struct{})); err == datastore.ErrNoSuchEntity {
		if !h.createIfMissing {
			http.Error(w, fmt.Sprintf("namespace %s not found", namespace), http.StatusNotFound)
			return
		}

		_, err := datastore.Put(ctx, key, new(struct{}))
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create namespace %s: %v", namespace, err), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, fmt.Sprintf("fetching namespace %s: %v", namespace, err), http.StatusInternalServerError)
		return
	}
	h.h(w, r, key)
}

func getAll(w http.ResponseWriter, r *http.Request, namespace *datastore.Key) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "getAll")

	values := []value{}
	keys, err := datastore.NewQuery(valueKind).Ancestor(namespace).GetAll(ctx, &values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	printValues := r.FormValue("v") == "true"
	if printValues {
		for i, key := range keys {
			fmt.Fprintf(w, "%s:%s\n", key.StringID(), values[i].Value)
		}
	} else {
		for _, key := range keys {
			fmt.Fprintln(w, key.StringID())
		}
	}
}

func getOne(w http.ResponseWriter, r *http.Request, namespace *datastore.Key) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "getOne")

	vars := mux.Vars(r)
	keyName := vars["key"]
	key := datastore.NewKey(ctx, valueKind, keyName, 0, namespace)

	var v value
	if err := datastore.Get(ctx, key, &v); err == datastore.ErrNoSuchEntity {
		http.Error(w, fmt.Sprintf("key %s not found in namespace %v", keyName, namespace.StringID()), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("fetching value %s in %s: %v", keyName, namespace.StringID(), err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintln(w, v.Value)
}

func put(w http.ResponseWriter, r *http.Request, namespace *datastore.Key) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "put")

	vars := mux.Vars(r)
	keyName := vars["key"]
	key := datastore.NewKey(ctx, valueKind, keyName, 0, namespace)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v := value{string(b)}

	if _, err := datastore.Put(ctx, key, &v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func delete(w http.ResponseWriter, r *http.Request, namespace *datastore.Key) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "delete")

	vars := mux.Vars(r)
	keyName := vars["key"]

	key := datastore.NewKey(ctx, valueKind, keyName, 0, namespace)
	if err := datastore.Delete(ctx, key); err != nil {
		http.Error(w, fmt.Sprintf("fetching value %s in %s: %v", keyName, namespace, err), http.StatusInternalServerError)
	}
}
