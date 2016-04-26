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
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

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

	// TODO: iterate over all the events and fetch the weather for its location.
	// If that operation fails just log the error and continue with the next event.
	// If it doesn't fail modify the weather in the event.

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
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
	// TODO: use apiTemplate above to fetch the weather for the location
	// You will need to get the APPID which is stored as an environment variable named API_KEY.
	// You need to replace the API_KEY defined in your app.yaml. Read the instructions there.
	// To access environment variables use the Getenv method from the os package.
	// Don't forget to URL escape the location before using it with the template.

	// TODO: use urlfetch to create a client and call the API.
	// Dont' forget to close the Body of the response at the end of this function.

	// TODO: Create a variable with anonymous structure type containing the fields we care about.
	// We care about the first weather of the list of weathers, including its description and icon.
	// We also need the error message to understand if something went wrong.
	// See the examples in good_api_output.json and bad_api_output.json to undertand
	// how the values are encoded.

	// TODO: Create a json decoder reading from the response's body
	// and extract the information we care about.

	// TODO: check wheter the error message is empty, if not return an error with its contents.

	// Return the first element of the list of weathers we decoded.

	return nil, errors.New("weather not implemented")
}
