package server

import (
	"net/http"
	"time"
)

type RouteMeta struct {
	Path    string
	Handler http.Handler
}

type option struct {
	routes []RouteMeta
}

type Option = func(*option)

func WithRoutes(m []RouteMeta) Option {
	return func(o *option) {
		o.routes = append(o.routes, m...)
	}
}

func New(options ...Option) *http.Server {
	mux := http.NewServeMux()

	o := &option{
		routes: make([]RouteMeta, 0),
	}

	for _, opt := range options {
		opt(o)
	}

	for _, r := range o.routes {
		mux.Handle(r.Path, r.Handler)
	}

	server := http.Server{
		Addr:        "localhost:1234",
		Handler:     mux,
		IdleTimeout: 30 * time.Second,
	}
	return &server
}
