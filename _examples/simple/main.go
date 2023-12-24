package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IroNEDR/router"
	"github.com/IroNEDR/router/middleware"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func helloName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := ctx.Value(router.PathParamsCtx{}).(router.Params)
	fmt.Fprintf(w, "Hello %s!", params["name"])
}

func main() {

	log := log.New(os.Stdout, "", log.LstdFlags)
	middleware.SetLogger(log)

	r := router.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/hello", http.HandlerFunc(helloWorld))
	r.Get("/hello/:name", http.HandlerFunc(helloName))

	log.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", r)
}
