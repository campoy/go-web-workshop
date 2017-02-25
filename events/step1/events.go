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

package events

import (
	"io"
	"net/http"
	"sync"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
)

// Event contains the information related to an event.
type Event struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}

var (
	// you can remove these values if you prefer
	// they're here just in case you implement listEvents before addEvent.
	events = []Event{
		{
			Title:       "Craft Conf",
			Description: "CRAFT is about software craftsmanship, which tools, methods, practices should be part of the toolbox of a modern developer and company.",
			Date:        time.Date(2016, 4, 26, 0, 0, 0, 0, time.Local),
			Location:    "Budapest",
		},
		{
			Title:       "Google I/O",
			Description: "Google I/O is for developers - the creative coders who are building what's next. Each year, we explore the latest in tech, mobile \u0026 beyond.",
			Date:        time.Date(2016, 5, 28, 0, 0, 0, 0, time.Local),
			Location:    "Mountain View",
		},
		{
			Title:       "GopherCon China",
			Description: "GOPHER'S BIGGEST PARTY",
			Date:        time.Date(2016, 5, 18, 0, 0, 0, 0, time.Local),
			Location:    "Beijing",
		},
	}
	mu sync.RWMutex
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/events", listEvents).Methods("GET")
	r.HandleFunc("/api/events", addEvent).Methods("POST")
	http.Handle("/", r)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Create a new context from the http request.
	// Grab a read lock from the global mutex and remember to release it at the end.
	// Set the header "Content-Type"" on the response to "application/json".
	// Create a json encoder that will write to the response writer.
	// And use it to encode the list of events.
	// Handle the error and log with log.Errorf if not nil.

	http.Error(w, "listEvents not implemented", http.StatusNotImplemented)
}

func addEvent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	e, err := decodeEvent(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof(ctx, "event decoded: %+v", e)

	// TODO: grab a write lock from the global mutex and release it at the end.
	// Add the decoded event to the global list of events.
	// Set the status code of the response to 201.
	http.Error(w, "addEvent not implemented", http.StatusNotImplemented)
}

func decodeEvent(r io.Reader) (*Event, error) {
	var data struct {
		Title       string
		Date        string
		Location    string
		Description string
	}

	// TODO: create a json decoder that will decode from the parameter io.Reader.
	// Use it to decode into data, and handle the error.

	// TODO: If the Title is empty return an error, you can use fmt.Errorf or errors.New.
	// If the Location is empty return an error.
	// Parse the Date using the format "2006-01-02", return an error if it is not valid.

	return &Event{
		Title: data.Title,
		// TODO: assign the rest of values here.
	}, nil
}
