# 1: Web Clients

In this section you're going to learn everything you need to know to write a web servers and clients in Go.

The `net/http` package provides a series of functions and types to help you sending HTTP requests.
The most important types are:

- the [Client](https://golang.org/pkg/net/http#Client)
- the [Request](https://golang.org/pkg/net/http#Request)
- the [Response](https://golang.org/pkg/net/http#Response)

We'll get to see how those types work in a minute.
But before that it is important to realize that there's helper functions that will make our life easier here too.

## The Get method

One of those helper functions is [`Get`](https://golang.org/pkg/net/http#Get).

[embedmd]:# (examples/get.go /package main/ $)
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// try changing the value of this url
	res, err := http.Get("https://golang.org")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Status)
}
```
[source code](examples/get.go)

This program will send a GET HTTP request to the Go homepage and will print the status code of the response,
unless something goes wrong in which case it will log the error and stop the execution of the program.

Try changing the value of the URL to see what other codes you're able to get.
Some ideas you could try:

- https://golang.org/foo
- https://thisurldoesntexist.com
- https:/thisisnotaurl

### Exercise status codes

The `Get` function we just used returns a `Response` and an `error`.
Read the documentation of `Response` and modify `get.go` so it prints a message only when the status code of the response is 404.

### Exercise body

The `Response` type also has a `Body` field of type `io.ReadCloser`.
The type of this field is a big hint on what we can do with it: read it and close it.

Modify the program `get.go` so it will also print the body of the `Response` onto the standard output.
Remember that the `Body` should be closed at the end of your program to avoid leaking memory!

## Other methods

The `Get` method we just used is a helper function that calls the `Get` method of the `http.DefaultClient`.
The `Client` type offers some other methods related directly to the HTTP methods we all know:

- [Client.Get](https://golang.org/pkg/net/http#Client.Get)
- [Client.Post](https://golang.org/pkg/net/http#Client.Post)
- [Client.PostForm](https://golang.org/pkg/net/http#Client.PostForm)
- [Client.Head](https://golang.org/pkg/net/http#Client.Head)

and the equivalent helper functions:

- [http.Get](https://golang.org/pkg/net/http#Get)
- [http.Post](https://golang.org/pkg/net/http#Post)
- [http.PostForm](https://golang.org/pkg/net/http#PostForm)
- [http.Head](https://golang.org/pkg/net/http#Head)

What if you want to use some other methods? Meet the [`Do`](https://golang.org/pkg/net/http#Client.Do) method.
This method receives a pointer to `http.Request` as parameters and return a `http.Response` and an `error`.

The `Request` method provides all the expressivity we need to send any kind of HTTP requests.

For instance, we can create the equivalent request to the original `get.go` program:

The `Client` type exposes the `Do` method that send the given `Request` and returns a `Response` and an `error`.

[embedmd]:# (examples/do-get.go /package main/ $)
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func doGet() {
	req, err := http.NewRequest("GET", "https://golang.org", nil)
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("http request failed: %v", err)
	}
	fmt.Println(res.Status)
}
```
[source code](examples/do-get.go)

This also allows us to provide a `Body` for the `Request`, similar to the one we had with the `Response`.
The big question is how to create an `io.Reader`?

Well, it depends on the type of what you want to read, but if you want to provide a `string` it is quite easy.
You can create an `io.Reader` from a `string` with [`strings.NewReader`](https://golang.org/pkg/strings#NewReader).

### Exercise PUT request

Modify the code above to send a PUT request to https://http-methods.appspot.com/YourName/Message.
Replace `YourName` with your name, or something unique that no one else might be using.
This will save whatever string you send in the Body of the request so you can retrieve it later with:

```
$ curl https://http-methods.appspot.com/YourName/Message
```

## Parameters: queries and forms

The web application behind https://http-methods.appspot.com supports printing all the keys under a namespace.
Try visiting https://http-methods.appspot.com/Hungary/ and you'll see all the keys in that namespace.
To also show the values you can add `?v=true` to the url: https://http-methods.appspot.com/Hungary/?v=true .

### Exercise queries

Write a program that will fetch and display all the keys and values in a namespace.
Rather than adding the `?v=true` at the end of the URL passed to `NewRequest` find the way to add it
directly in the `Request` that you pass to `Do`.

# Congratulations!

Good job! You now know pretty much everything there is to know (today) for HTTP clients, requests, and responses.
Now, to understand how all of this is handled on the server side, let's go the [next section](../section02/README.md).
