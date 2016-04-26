// Copygright 2016 Google Inc. All rights reserved.
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
	"net/url"
	"os"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
)

const eventKind = "Event"

// Event contains the information related to an event.
type Event struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Weather     Weather   `json:"weather" datastore:"-"`
}

// Weather contains the description and icon for a weather condition.
type Weather struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/events", listEvents).Methods(http.MethodGet)
	r.HandleFunc("/api/events", addEvent).Methods(http.MethodPost)
	http.Handle("/", r)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	events := []Event{}
	q := datastore.NewQuery(eventKind).
		Filter("Date >", time.Now()).
		Order("Date").
		Limit(5)

	_, err := q.GetAll(ctx, &events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, e := range events {
		w, err := weather(ctx, e.Location)
		if err != nil {
			log.Errorf(ctx, "fetching weather for %q: %v", e.Location, err)
			continue
		}
		events[i].Weather = *w
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Errorf(ctx, "encoding events: %v", err)
	}
}

func addEvent(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("BLOCK_WRITES") == "True" {
		http.Error(w, "this is a read only instance, sorry", http.StatusForbidden)
		return
	}

	ctx := appengine.NewContext(r)

	e, err := decodeEvent(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := datastore.NewIncompleteKey(ctx, eventKind, nil)
	_, err = datastore.Put(ctx, key, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodeEvent(r io.Reader) (*Event, error) {
	var data struct {
		Title       string
		Date        string
		Location    string
		Description string
	}

	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode json: %v", err)
	}

	if data.Title == "" {
		return nil, fmt.Errorf("title can't be empty")
	}
	if data.Location == "" {
		return nil, fmt.Errorf("location is required")
	}
	t, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return nil, fmt.Errorf("parse date: %v", err)
	}

	return &Event{
		Title:       data.Title,
		Date:        t,
		Description: data.Description,
		Location:    data.Location,
	}, nil
}

const (
	apiTemplate     = "http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s"
	iconURLTemplate = "http://openweathermap.org/img/w/%s.png"
)

func weather(ctx context.Context, location string) (*Weather, error) {
	// check if the weather for the location is in memcache.
	var weather Weather
	_, err := memcache.JSON.Get(ctx, location, &weather)
	if err == nil {
		return &weather, nil
	} else if err != memcache.ErrCacheMiss {
		log.Errorf(ctx, "could not retrieve weather for %q from memcache: %v", location, err)
	}

	// prepare the request to the weather API.
	apiKey := os.Getenv("API_KEY")
	location = url.QueryEscape(location)
	api := fmt.Sprintf(apiTemplate, apiKey, location)

	client := urlfetch.Client(ctx)
	res, err := client.Get(api)
	if err != nil {
		return nil, fmt.Errorf("could not get weather: %v", err)
	}

	// we need to close the body of the API response to avoid leaks.
	defer res.Body.Close()

	// we need to decode the list of weathers and the error message.
	var data struct {
		Weather []Weather
		Message string
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("could not decode weather: %v", err)
	}

	// if the error message is not empty, something bad happened.
	if data.Message != "" {
		return nil, fmt.Errorf("no weather found: %s", data.Message)
	}

	// we just take the first value for the weather.
	weather = data.Weather[0]
	// and make the icon a complete url.
	weather.Icon = fmt.Sprintf(iconURLTemplate, weather.Icon)

	// cache the weather in memcache for later.
	item := &memcache.Item{
		Key:        location,
		Object:     &weather,
		Expiration: 1 * time.Hour,
	}
	err = memcache.JSON.Set(ctx, item)
	if err != nil {
		log.Errorf(ctx, "could not cache weather for %q: %v", location, err)
	}

	return &weather, nil
}
