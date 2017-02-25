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
	"errors"

	"golang.org/x/net/context"
)

const (
	apiURL          = "http://api.openweathermap.org/data/2.5/weather"
	iconURLTemplate = "http://openweathermap.org/img/w/%s.png"
)

func weather(ctx context.Context, location string) (*Weather, error) {
	// TODO: use apiURL above to fetch the weather for the location
	// You will need to create a map of URL parameters, known as url.Values from the "net/url" package.
	// Set the parameter "APPID" to the environment variable WEATHER_API_KEY defined in your app.yaml.
	// To access environment variables use the Getenv method from the os package.
	// Then set the parameter "q" to the location.

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
