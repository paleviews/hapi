package runtime

import (
	"net/http"
)

type Service struct {
	RPCs []RPC
}

type RPC struct {
	Route   Route
	Handler func(http.ResponseWriter, *http.Request)
}

type Route struct {
	Method HTTPMethod
	Path   string
}

type HTTPMethod string

const (
	HTTPMethodGet    HTTPMethod = "GET"
	HTTPMethodPost   HTTPMethod = "POST"
	HTTPMethodDelete HTTPMethod = "DELETE"
	HTTPMethodPut    HTTPMethod = "PUT"
	HTTPMethodPatch  HTTPMethod = "PATCH"
)
