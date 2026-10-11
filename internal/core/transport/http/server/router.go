package core_http_server

import (
	"net/http"

	core_http_middleware "todoapp/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	ApiVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	routes     []Route
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion APIVersion,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) AddRoutes(routes ...Route) {
	r.routes = append(r.routes, routes...)
}

func (r *APIVersionRouter) Handlers() map[string]http.Handler {
	handlers := make(map[string]http.Handler, len(r.routes))

	for _, route := range r.routes {
		pattern := route.Method + " /api/" + string(r.apiVersion) + route.Path
		handler := core_http_middleware.ChainMiddleware(
			route.WithMiddleware(),
			r.middleware...,
		)

		handlers[pattern] = handler
	}

	return handlers
}
