# 3: Input validation and status codes

So far we've just assumed that all the HTTP requests our server receives are good.
And if there's something that you should never do when writing web servers is trusting your input!

In this chapter we'll see how to extract the information sent in different parts of the HTTP request.
Once we have that information we'll see how can validate them, and how we can signal different errors.

Let's start!

## Reading parameters from the URL

We've seen how we can route a request to different handlers depending on the path.
Now we're going to see how to extract then ones in the query part of the request, aka the data after `?`.

The `http.Request` type has a method `FormValue` with the following docs:

    func (r *Request) FormValue(key string) string

    FormValue returns the first value for the named component of the query. POST and PUT body parameters take precedence over URL query string values. FormValue calls ParseMultipartForm and ParseForm if necessary and ignores any errors returned by these functions. If key is not present, FormValue returns the empty string. To access multiple values of the same key, call ParseForm and then inspect Request.Form directly.

That's easy! So if we want to obtain the value of a parameter `q` in the URL `/hello?msg=world`
we can write the next program.

[embedmd]:# (examples/handlers/main.go /func paramHandler/ /^}/)
```go
func paramHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		name = "friend"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}
```

### Exercise Hello, parameter

Write a web server that will answer to requests to `/hello?name=world` with an HTTP response with the text `Hello, world!`.
If the `name` is not present it should print `Hello, friend!`.

You can test it with your own browser, but let's try a couple things with `curl` too.
Before running these think about what you expect to see and why.

```bash
$ curl "localhost:8080/hello?name=world"

$ curl "localhost:8080/hello?name=world&name=francesc"

$ curl -X POST -d "name=francesc" "localhost:8080/hello"

$ curl -X POST -d "name=francesc" "localhost:8080/hello?name=potato"
```

Think about how would you make your program print all the values given to `name`.

## Reading from the Request body

Similarly to how we read from the `Body` in the `http.Response` a couple of chapter before we can read
the body of the `http.Request`.

Note that even though the type of `Body` in `http.Request` is `io.ReadCloser` the body will be automatically
closed at the end of the execution of the http handler, so don't worry about it.

There's many ways we can read from an `io.Reader`, but for now you can use `ioutil.ReadAll`,
which returns a `[]byte` and an `error` if something goes wrong.

[embedmd]:# (examples/handlers/main.go /func bodyHandler/ /^}/)
```go
func bodyHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "could not read body: %v", err)
		return
	}
	name := string(b)
	if name == "" {
		name = "friend"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}
```

### Exercise Hello, body

Modify the previous exercise so instead of reading the `name` argument from a query or form value it will
use whatever content of the `Body` is. Again, if the `Body` is empty the response should greet `friend`.

If the call to `ioutil.ReadAll` returns an error write that error message to the output.

You can test your exercise by using `curl`:

```bash
$ curl -X POST -d "francesc" "localhost:8080/hello"
```

As an extra exercise remove any extra blank spaces surrounding the name.

## Communicating errors

In the previous exercise we decided to just print the error if the call to `ioutil.ReadAll` failed.
That is actually a pretty horrible idea, as you might imagine 😁.

How should we do it? Well, the HTTP protocol defines a set of status codes that help us describe the nature of a response.
We've actually used them before, when we checked if the response obtained using `Get` was `OK` or not.

There are two ways of setting the status code with a `ResponseWriter`.

### Status codes with ResponseWriter.WriteHeader

Using the `WriteHeader` method in `ResponseWriter` we can set the status code of the response.
The parameter is an `int` so you could pass any number, but it is better to use the constants already
defined in the `http` package. They all start with `Status` and you can find them [here](https://golang.org/pkg/net/http/#pkg-constants).

By default, the status code of a response will be `StatusOK` aka `200`.

#### Exercise better errors

Modify your previous program so the response in case of error will have status code `500`.
Instead of using the number `500` find the corresponding constant.

Then make the status of the response be `400` when the body is empty.

### Status codes with http.Error

In the previous exercise you've noticed that often when we set the status code of the response we also
write the description of the error. That's why the `http.Error` function exists.

    func Error(w ResponseWriter, error string, code int)

    Error replies to the request with the specified error message and HTTP code. The error message should be plain text.

A call to `Error` can then replace a call to `WriteHeader` followed by a some call writing to the `ResponseWriter`.

#### Exercise status codes with http.Error

Replace your calls to `WriteHeader` and `Fprintf` in the previous exercise with a call to `Error`.

### Response headers

You might have noticed that if you send more than one line in the response,
your browser shows it as one. Why is that?

The answer is that the `net/http` packages guesses that the output is HTML,
and therefore your browser concatenates the lines. For short outputs, it's hard to guess.
You can see the content type of your response by adding `-v` to your `curl` command.

```bash
$ curl -v localhost:8080
< HTTP/1.1 200 OK
< Date: Mon, 25 Apr 2016 16:14:46 GMT
< Content-Length: 19
< Content-Type: text/html; charset=utf-8
```

So, how do we stop the `net/http` package from guessing the content type? We specify it!
To do so we need to set the header `"Content-Type"` to the value `"text/plain"`.
You can set headers in the response with the `Header` function in the `ResponseWriter`.

`Header` returns a [`http.Header`](https://golang.org/pkg/net/http/#Header) which has, among other methods,
the method `Set`. We can then set the content type in our `ResponseWriter` named `w` like this.

[embedmd]:# (examples/texthandler/main.go /w.Header.*/)
```go
w.Header().Set("Content-Type", "text/plain")
```

### Avoiding repetition

Imagine if you had hundreds of different http handlers, and you want to set the content type header
in each and every one of them. That sounds painful, doesn't it?

Let me tell you about a cool technique that helps defining behavior shared by many handlers.
Some people call them decorators, most of them also write Python 😛.

To start we're going to define a new type named `textHandler` that contains a `http.HandlerFunc`.

[embedmd]:# (examples/texthandler/main.go /type textHandler/ /^}/)
```go
type textHandler struct {
	h http.HandlerFunc
}
```

Now we're going to define the `ServeHTTP` method on `textHandler` so it satisfies the `http.Handler` interface.

[embedmd]:# (examples/texthandler/main.go /.*ServeHTTP/ /^}/)
```go
func (t textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set the content type
	w.Header().Set("Content-Type", "text/plain")
	// Then call ServeHTTP in the decorated handler.
	t.h(w, r)
}
```

Finally we replace our `http.HandleFunc` calls with `http.Handle`.

[embedmd]:# (examples/texthandler/main.go /func main/ /^}/)
```go
func main() {
	http.Handle("/hello", textHandler{helloHandler})
	http.ListenAndServe(":8080", nil)
}
```

#### Exercise Setting headers only once

Modify the program from your previous exercise so you have only one line that set the content type of the response.

#### Exercise Even better error handling (optional)

Modify the `textHandler` that we showed before so instead of `http.HandlerFunc` it receives a function that
returns an `int` and an `error`. The `ServeHTTP` method of `textHandler` should check use that integer and
the `error` and set the status code and content accordingly.

_Super mega optional_: define a new error type that contains the information about the status code too.
The first to finish this might get a prize ... just saying.

# Congratulations!

You're now able to validate the input to your http handlers and set the status code and content type accordingly.
You are clearly awesome! 🎉

But you can be awesomer, I assure you, by going to [section 4](../section04/README.md).
