/*
 * Copyright 2017 Google Inc. All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to writing, software distributed
 * under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.
 *
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

function EventsCtrl($scope, $http) {
  // The list of events in display.
  $scope.events = [];
  // The fields in the event creation dialog.
  $scope.newEvent = {};

  // Display an error using an alert dialog.
  var alertError = function(data, status) {
    alert('code ' + status + ': ' + data);
  };

  // Fetches all the events from the API.
  var fetchEvents = function() {
    return $http.get('/api/events').
      error(alertError).
      success(function(data) { $scope.events = data; });
  };

  // Adds a new event throught the API.
  $scope.addEvent = function() {
    $http.post('/api/events', $scope.newEvent).
      error(alertError).
      success(function() {
        fetchEvents().then(function () {
          // If everything worked, clear the dialog.
          $scope.event = {};
          // Fetch again after a bit, in case of eventual consistency.
          setTimeout(fetchEvents, 1000);
        });
      });
  };

  // Fetch the list of events from the API.
  fetchEvents();
}
