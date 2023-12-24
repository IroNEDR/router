package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IroNEDR/router"
	"github.com/IroNEDR/router/middleware"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Try out the following endpoints:\n\tGET /category/:category/:product\n\tPOST /category/:category/:product?color=color")
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := ctx.Value(router.PathParamsCtx{}).(router.Params)
	queryParams := ctx.Value(router.QueryParamsCtx{}).(router.Params)
	fmt.Fprintf(w, "Overview of product %s in %s from category %s",
		params["product"], params["category"], queryParams["color"])
}

func orderProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := ctx.Value(router.PathParamsCtx{}).(router.Params)
	queryParams := ctx.Value(router.QueryParamsCtx{}).(router.Params)
	if queryParams["color"] == "" {
		queryParams["color"] = "white"
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Order Successful:\n\tproduct %s in %s from category %s",
		params["product"], params["category"], queryParams["color"])
}

func main() {

	log := log.New(os.Stdout, "", log.LstdFlags)
	middleware.SetLogger(log)

	r := router.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", http.HandlerFunc(welcome))
	r.Get("/category/:category/:product", http.HandlerFunc(getProduct))
	r.Post("/category/:category/:product", http.HandlerFunc(orderProduct))

	log.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", r)
}
