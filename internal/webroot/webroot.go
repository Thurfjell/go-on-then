package webroot

import (
	"embed"
	"net/http"
)

type RouteMeta struct {
	Path    string
	Handler http.Handler
}

type root struct {
	routes []RouteMeta
}

type Option = func(*root)

func WithRouterMeta(meta []RouteMeta) Option {
	return func(r *root) {
		r.routes = append(r.routes, meta...)
	}
}

//go:embed all:static/*
var static embed.FS

func New(options ...Option) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	rootMeta := &root{
		routes: make([]RouteMeta, 0),
	}

	for _, opt := range options {
		opt(rootMeta)
	}

	for _, r := range rootMeta.routes {
		mux.Handle(r.Path, r.Handler)
	}

	return mux
}
