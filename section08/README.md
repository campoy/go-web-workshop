# 8: Retrieving remote resources with urlfetch

Sometimes your application will need to communicate with the external world,
send data via POST requests or maybe retrieve some information using GET.

The App Engine framework limits what you can do to ensure scalability and
performance, which means you can't use the `net/http` package directly to
run `http.Get("https://google.com")` but it's just as simple using the
[`appengine/urlfetch`](https://cloud.google.com/appengine/docs/go/urlfetch/)
package provided by the App Engine runtime:

The most important function of the `urlfetch` package is `Client`:

```go
func Client(context appengine.Context) *http.Client
```

So given an `appengine.Context` you get an HTTP client, and then you can start
from there. So if you wanted to fetch Google's home page you would do:

[embedmd]:# (fetch/fetch.go /package fetch/ /^}/)
```go
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
```

If you run this code which you can find in this [directory](fetch) you should
see a slightly broken version of the Google homepage, as all the local links
will be broken.

# Exercise: fetching weather for event locations

With what you just learned and your previous knowledge on JSON encoding and
decoding, it is time to tackle [step 3](../events/step3/README.md).

# Congratulations!

You are now capable of making your App Engine application communicate with the
rest of the web using HTTP requests thanks to the `urlfetch` package!

But is it really a good idea to fetch external resources every time we need to
serve a request? Shouldn't we be caching some of this data for later use?

Well, continue to the [next chapter](../section09/README.md) to learn how to
cache and retrieve information from App Engine using memcache.
