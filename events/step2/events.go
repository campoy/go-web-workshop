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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
)

const eventKind = "Event"

// Event contains the information related to an event.
type Event struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/events", listEvents).Methods("GET")
	r.HandleFunc("/api/events", addEvent).Methods("POST")
	http.Handle("/", r)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var events []Event

	// TODO: Create a new context from the http request
	// Create a new query on the kind event.
	// Filter all the events that are in the past.
	// Order the events by date.
	// And limit so only 5 events are shown.

	// TODO: Use the GetAll method to fetch all the events into the slice of Events.
	// Check for any errors.

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Errorf(ctx, "encoding events: %v", err)
	}
}

func addEvent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	e, err := decodeEvent(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Infof(ctx, "event decoded: %v", e)

	// TODO: Create a new incomplete key of type Event.
	// And use it to store the decoded event.
	// Remember to handle the error!
	http.Error(w, "add event not implemented", http.StatusNotImplemented)
	return

	w.WriteHeader(http.StatusCreated)
}

func decodeEvent(r io.Reader) (*Event, error) {
	// Using an anonymous type instead of Event because the date format
	// we need to parse is not standard.
	// You can find more on this in this talk: https://talks.golang.org/2015/json.slide
	var data struct {
		Title       string
		Date        string
		Location    string
		Description string
	}

	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("could not decode JSON: %v", err)
	}

	if data.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if data.Location == "" {
		return nil, fmt.Errorf("location is required")
	}
	t, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}

	return &Event{
		Title:       data.Title,
		Date:        t,
		Description: data.Description,
		Location:    data.Location,
	}, nil
}
