# 2: Web servers

In this section you'll learn how to write a simple HTTP server in Go.

We will use the [`net/http`](https://golang.org/pkg/net/http) package to do so, click
on that link to browse its documentation.

### Defining HTTP handlers

The `net/http` package defines the [`HandlerFunc`](https://golang.org/pkg/net/http#HandlerFunc) type:

```go
type HandlerFunc func(ResponseWriter, *Request)
```

The first parameter of this function type is a
[`ResponseWriter`](https://golang.org/pkg/net/http#ResponseWriter), which
provides a way to set headers on the HTTP response. It also provides a `Write`
method which makes it satisfy the `io.Writer` interface.

Let's see a very simple HTTP handler that simply writes `"Hello, web"` to the output:

[embedmd]:# (examples/step1/main.go /package main/ /^}/)
```go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, web")
}
```

As you can see we're using the `fmt.Fprintln` function, whose first parameter
is an `io.Writer`.

### Registering HTTP handlers

Once a handler is defined we need to inform the `http` package about it and
specify when to run it. To do so we can use the `http.HandleFunc` function:

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

The first parameter is a pattern which will be used to decide when to execute a
handler, and the second argument is the handler itself.

Patterns name fixed, rooted paths, like `"/favicon.ico"`, or rooted subtrees,
like `"/images/"` (note the trailing slash). Longer patterns take precedence
over shorter ones, so that if there are handlers registered for both
`"/images/"` and `"/images/thumbnails/"`, the latter handler will be called for
paths beginning `"/images/thumbnails/"` and the former will receive requests
for any other paths in the `"/images/"` subtree.

Let's see how to register our `helloHandler` defined above:

[embedmd]:# (examples/step2/main.go /package main/ $)
```go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, web")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
}
```

Note that we're registering our handler as part of the `main` function.

Try to run the code above:

```bash
$ go run examples/step2/main.go
```

What happens? Well, we're missing the last piece of the puzzle: starting the
web server!

#### The Handler interface

Using `http.HandleFunc` and passing a value of type `http.HandlerFunc` can be pretty constraining.
There's also another function `http.Handle` that will accept any value satisfying the `http.Handler` interface.

The [`http.Handler`](https://golang.org/pkg/net/http/#Handler) interface is defined in the `http` package as:

```go
type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
}
```

And guess what, the type `http.HandlerFunc` satisfies `http.Handler` thanks to
[`HandlerFunc.ServeHTTP`](https://golang.org/pkg/net/http/#HandlerFunc.ServeHTTP).

The code of the `ServeHTTP` method for `HandlerFunc` is something of beauty.

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

We will see how this interface is the extension points where web frameworks and toolkits add functionality.

### Listening and serving

Once the handlers have been defined and registered we need to start the HTTP
server to listen for requests and execute the corresponding handler.

To do so we use the function
[`http.ListenAndServe`](https://golang.org/pkg/net/http#ListenAndServe):

```go
func ListenAndServe(addr string, handler Handler) error
```

The first parameter is the address on which we want the server to listen,
we could use something like `"127.0.0.1:8080"` or `"localhost:80"`.

The second parameter is an `http.Handler`, a type that allows you to define
different ways of handling requests. Since we're using the default methods
with `HandleFunc` we don't need to provide any value here: `nil` will do.

And last but definitely not least the function returns an `error`. In Go,
errors are handled by returning values rather than throwing exceptions.

The type `error` is a predefined type (just like `int` or `bool`) and is an interface
with only one method:

```go
type error interface {
	Error() string
}
```

By convention errors are the last value returned by methods and functions and
when no error has occurred the returned value equals to `nil`.

So if we want to check that our server started successfully and log an error
otherwise we would modify our code to add a call to `ListenAndServe`.

[embedmd]:# (examples/step3/main.go /package main/ $)
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, web")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

Running this code should now start a web server listening on `127.0.0.1:8080`.

Try it:

	$ go run main.go

And then visit http://127.0.0.1:8080/hello.

### Exercise Bye, web

Modify the program above by adding a second handler named `byeHandler` that prints `"Bye, web"`
to the http response.

### Exercise Hello, Handler

Modify the program from the previous example so you can replace the call to `http.HandleFunc`
by a call to `http.Handle`. You will need to define a new type `helloHandler` and make that type
satisfy the `http.Handler` interface.

### A better multiplexor

Soon you will start having more complicated requirements to route your requests
to your handlers such as:

- route depending on the methods: `POST` and `GET` routed different handlers.
- variable extraction from paths: `/product/{productID}/part/{partID}`

These cases can be handled either by hand or using a toolkit that will plug
correctly into the existing `net/http` package, such as the
[Gorilla toolkit](http://www.gorillatoolkit.org/) and its `mux` package.

[embedmd]:# (examples/gorilla.go /package main/ $)
```go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func listProducts(w http.ResponseWriter, r *http.Request) {
	// list all products
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	// add a product
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["productID"]
	log.Printf("fetching product with ID %q", id)
	// get a specific product
}

func main() {
	r := mux.NewRouter()
	// match only GET requests on /product/
	r.HandleFunc("/product/", listProducts).Methods("GET")

	// match only POST requests on /product/
	r.HandleFunc("/product/", addProduct).Methods("POST")

	// match GET regardless of productID
	r.HandleFunc("/product/{productID}", getProduct)

	// handle all requests with the Gorilla router.
	http.Handle("/", r)
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}
```

Gorilla also provides packages for session management, cookies, etc.
Have a look at the [documentation](www.gorillatoolkit.org/).

#### Exercise Hello, {you}

Using the `mux` package from the previous example write a new web server.
This server will handle all HTTP requests sent to `/hello/name` with an HTTP page
containing the text `"Hello, name"`. The `name` in this example can of course change,
so if the request was `/hello/Francesc` the response should say `"Hello, Francesc"`.

_Note_: to install the `mux` package in your machine you can use `go get`:

```bash
$ go get github.com/gorilla/mux
```

# Congratulations!

You just wrote your first HTTP server in Go! Isn't it awesome? Well, it doesn't
do much yet but the best is to come.

On the next chapter we'll learn how to validate whatever the input of your HTTP endpoints and
how to signal different problems in the HTTP responses.

Continue to [the next section](../section03/README.md).
