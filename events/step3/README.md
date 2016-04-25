# Step 3: adding weather with openweathermap.org

So far, we have an application that is capable of storing events submitted
through a web form, and display them back into the web page.

Now we're going to add some cool functionality by retrieving the weather for
the location where the event is taking place. To do so we need to fetch
information from an external API, and we will use open openweathermap.org.

You need to sign up to https://openweathermap.org to obtain an API key and replace the
value of WEATHER_API_KEY in the app.yaml file.

Then you'll implement the call to the API, preparing the request, decoding the
response, and adding the resulting weather to the events.

This will cause one API call per event per request, which is incredibly wasteful.
But don't worry we'll fix that later, for now come back to the
[instructions](../../section08/README.md#congratulations).
