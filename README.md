# Router

[![Go Report Card](https://goreportcard.com/badge/github.com/IroNEDR/router)](https://goreportcard.com/report/github.com/IroNEDR/router)

This minimalist Go router is designed for learning purposes to help you understand how to create your own router and implement supporting middleware in Go without using any third-party libraries. It provides a simple yet functional routing mechanism with the ability to add middleware functions for enhanced request handling.

## Features

- Custom routing based on HTTP methods and path patterns.
- Support for parameterized paths.
- Middleware support for request processing.
- A basic `NotFoundHandler` for handling 404 errors.

## Installation

To use this router in your Go project, run:
`go get github.com/IroNEDR/router`

## Example

```go
package main

import (
    "fmt"
	"log"
	"net/http"
	"os"

	"github.com/IroNEDR/router"
	"github.com/IroNEDR/router/middleware"
)

// calling http://localhost:8080/hello/John?age=25 will print "Hello John, age 25!"
func helloName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := ctx.Value(router.PathParamsCtx{}).(router.Params)
    queryParams := ctx.Value(router.QueryParamsCtx{}).(router.Params)
    if queryParams["age"] == "" {
        queryParams["age"] = "unknown"
    }
    w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello %s!", params["name"], queryParams["age"])
}

func main() {
    
	log := log.New(os.Stdout, "", log.LstdFlags)
	middleware.SetLogger(log)
	r := router.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/hello/:name", http.HandlerFunc(helloName))
    log.Println("Server listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
```
More examples can be found in the `_examples` directory.
