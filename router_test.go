package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutePathMatch(t *testing.T) {
	rp := routePath("/hello/:name")
	params := rp.match("/hello/John")

	if params == nil {
		t.Error("params should not be nil")
	}
	if params["name"] != "John" {
		t.Errorf("params[\"name\"] should be \"John\", got %s", params["name"])
	}

	params = rp.match("/order/123")
	if params != nil {
		t.Error("params should be nil")
	}
}

func TestExtractQueryParams(t *testing.T) {
	params := extractQueryParams("name=John&age=30")
	if params["name"] != "John" && params["age"] != "30" {
		t.Errorf("params should be {\"name\": \"John\", \"age\": \"30\"}, got %s", params)
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()

	req, err := http.NewRequest("GET", "/hello/John", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := ctx.Value(PathParamsCtx{}).(Params)
		if params["name"] != "John" {
			t.Errorf("params[\"name\"] should be \"John\", got %s", params["name"])
		}
	})

	r.Get("/hello/:name", testHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}
