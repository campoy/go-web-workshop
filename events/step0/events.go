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

func init() {
	// TODO: Create a new mux.Router
	// All the GET requests to /api/events should go to the handler listEvents.
	// All the POST requests to /api/events should go to the handler addEvent.
	// And all requests will be handled with the mux.Router.
}

// TODO: define a http.HandleFunc named listEvents
// When called it should simply print the contents of the listOutput constant declared below.
// You can also create a new App Engine context and use it to log some message.

// TODO: define a http.HandleFunc named addEvent
// When called it should simply set the status code to 201 and log a message.

const listOutput = `
[
    {
        "title": "Craft Conf",
        "description": "CRAFT is about software craftsmanship, which tools, methods, practices should be part of the toolbox of a modern developer and company.",
        "date": "2016-04-26T00:00:00Z",
        "location": "Budapest"
    },
    {
        "title": "Google I/O",
        "description": "Google I/O is for developers - the creative coders who are building what's next. Each year, we explore the latest in tech, mobile \u0026 beyond.",
        "date": "2016-04-28T00:00:00Z",
        "location": "Mountain View"
    },
    {
        "title": "GopherCon China",
        "description": "GOPHER'S BIGGEST PARTY",
        "date": "2016-05-18T00:00:00Z",
        "location": "Beijing"
    }
]
`
