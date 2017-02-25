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

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func listProducts(w http.ResponseWriter, r *http.Request) {
	// list all products
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	// add a product
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["productID"]
	log.Printf("fetching product with ID %q", id)
	// get a specific product
}

func main() {
	r := mux.NewRouter()
	// match only GET requests on /product/
	r.HandleFunc("/product/", listProducts).Methods("GET")

	// match only POST requests on /product/
	r.HandleFunc("/product/", addProduct).Methods("POST")

	// match GET regardless of productID
	r.HandleFunc("/product/{productID}", getProduct)

	// handle all requests with the Gorilla router.
	http.Handle("/", r)
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}
