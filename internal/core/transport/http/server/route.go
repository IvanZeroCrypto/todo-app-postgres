package core_http_server

import (
	"net/http"

	core_http_middleware "gitverse.ru/IvanVerse/go-todoApp/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

func (r *Route) WithMiddleWare() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}

// func NewRoute(method string, path string, handler http.HandlerFunc) Route {
// 	return Route{
// 		Method:  method,
// 		Path:    path,
// 		Handler: handler,
// 	}
// }
