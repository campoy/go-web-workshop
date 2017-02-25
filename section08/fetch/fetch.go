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

package fetch

import (
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// create a new HTTP client
	c := urlfetch.Client(ctx)

	// and use it to request the Google homepage
	res, err := c.Get("https://google.com")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// we need to close the body at the end of this function
	defer res.Body.Close()

	// then we can dump the whole webpage onto our output
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Errorf(ctx, "could not copy the response: %v", err)
	}
}

func init() {
	http.HandleFunc("/", handler)
}
