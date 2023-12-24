package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/IroNEDR/router/middleware"
)

type routePath string

type route struct {
	httpMethod string
	path       routePath
	handler    http.Handler
}

// Params is a multipurpose map for path and query parameters.
type Params map[string]string

// PathParamsCtx is the context key for path parameters.
type PathParamsCtx struct{}

// QueryParamsCtx is the context key for query parameters.
type QueryParamsCtx struct{}

// Router is a simple HTTP router that matches the request path and method
// and calls the corresponding handler. If no route matches, the NotFoundHandler
// is called.
type Router struct {
	routes          []route
	NotFoundHandler http.Handler
	Middlewares     []func(http.Handler) http.Handler
}

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{
		NotFoundHandler: middleware.Logger(http.HandlerFunc(NotFoundHandler)),
	}
}

func (rp *routePath) match(path string) Params {

	patternPath := strings.Split(string(*rp), "/")
	requestPath := strings.Split(path, "/")

	if len(patternPath) != len(requestPath) {
		return nil
	}

	params := make(Params)
	for i, p := range strings.Split(string(*rp), "/") {
		if len(p) > 0 && p[0] == ':' {
			params[p[1:]] = requestPath[i]
		} else if p != requestPath[i] {
			return nil
		}
	}
	return params
}

func extractQueryParams(path string) Params {
	params := make(Params)
	for _, p := range strings.Split(path, "&") {
		if len(p) > 0 {
			pair := strings.Split(p, "=")
			params[pair[0]] = pair[1]
		}
	}
	return params
}

// NotFoundHandler is the default handler for when no route matches.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func chain(middlewares []func(http.Handler) http.Handler, handler http.Handler) http.Handler {
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	return handler
}

// Use adds a middleware to the router. Middlewares are called in the order they are added.
// The middleware is only applied to routes added after the middleware is added.
func (rt *Router) Use(middleware func(http.Handler) http.Handler) {
	rt.Middlewares = append(rt.Middlewares, middleware)
}

func (rt *Router) addRoute(method, path string, handler http.Handler) {
	handler = chain(rt.Middlewares, handler)
	route := route{httpMethod: method, path: routePath(path), handler: handler}
	rt.routes = append(rt.routes, route)
}

// Get adds a GET route to the router.
func (rt *Router) Get(path string, handler http.Handler) {
	rt.addRoute(http.MethodGet, path, handler)
}

// Post adds a POST route to the router.
func (rt *Router) Post(path string, handler http.Handler) {
	rt.addRoute(http.MethodPost, path, handler)
}

// Put adds a PUT route to the router.
func (rt *Router) Put(path string, handler http.Handler) {
	rt.addRoute(http.MethodPut, path, handler)
}

// Delete adds a DELETE route to the router.
func (rt *Router) Delete(path string, handler http.Handler) {
	rt.addRoute(http.MethodDelete, path, handler)
}

// Patch adds a PATCH route to the router.
func (rt *Router) Patch(path string, handler http.Handler) {
	rt.addRoute(http.MethodPatch, path, handler)
}

// ServeHTTP implements the http.Handler interface, making the router
// compatible with the standard library.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rt.routes {
		params := route.path.match(r.URL.Path)
		if params != nil && route.httpMethod == r.Method {

			ctx := context.WithValue(r.Context(), PathParamsCtx{}, params)

			if r.URL.RawQuery != "" {
				queryParams := extractQueryParams(r.URL.RawQuery)
				ctx = context.WithValue(ctx, QueryParamsCtx{}, queryParams)
			}

			route.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}
	rt.NotFoundHandler.ServeHTTP(w, r)
}
