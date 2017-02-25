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
	"google.golang.org/appengine/datastore"
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
	Weather     *Weather  `json:"weather" datastore:"-"`
}

// Weather contains the description and icon for a weather condition.
type Weather struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/events", listEvents).Methods("GET")
	r.HandleFunc("/api/events", addEvent).Methods("POST")
	http.Handle("/", r)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	events := []Event{}
	q := datastore.NewQuery(eventKind).
		Filter("Date >", time.Now()).
		Order("Date").
		Limit(5)

	if _, err := q.GetAll(ctx, &events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, e := range events {
		w, err := weather(ctx, e.Location)
		if err != nil {
			log.Errorf(ctx, "fetching weather for %q: %v", e.Location, err)
			continue
		}
		events[i].Weather = w
	}

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

	key := datastore.NewIncompleteKey(ctx, eventKind, nil)
	if _, err := datastore.Put(ctx, key, e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
