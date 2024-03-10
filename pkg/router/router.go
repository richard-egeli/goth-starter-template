package router

import (
	"net/http"
	"strings"
)

type Router struct {
	root        string
	parent      *Router
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func New() *Router {
	return &Router{
		root:        "",
		mux:         http.NewServeMux(),
		middlewares: nil,
	}
}

func reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (router *Router) getMiddlewares() []func(http.Handler) http.Handler {
	if router.parent != nil {
		return append(router.parent.getMiddlewares(), router.middlewares...)
	} else {
		return router.middlewares
	}
}

func (router *Router) set(path string, method Method, handlers []func(http.Handler) http.Handler) {
	path = router.root + path
	length := len(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if len(method) > 0 {
		path = string(method) + " " + path
	}

	var handler http.Handler
	for _, h := range append(reverse(handlers), reverse(router.getMiddlewares())...) {
		handler = h(handler)
	}

	if strings.HasSuffix(path, "/") && length > 1 {
		absPath := strings.TrimSuffix(path, "/")
		router.mux.Handle(absPath, handler)
	}

	router.mux.Handle(path, handler)
}

func (router *Router) route() string {
	if router.parent != nil {
		return router.parent.route() + router.root
	}

	return router.root
}

func (router *Router) Use(middleware func(http.Handler) http.Handler) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *Router) Sub(subroute string) *Router {
	return &Router{
		root:        router.route() + subroute,
		mux:         router.mux,
		parent:      router,
		middlewares: nil,
	}
}

func (router *Router) Get(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, GET, handlers)
}

func (router *Router) Post(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, POST, handlers)
}

func (router *Router) Delete(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, DELETE, handlers)
}

func (router *Router) Put(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, PUT, handlers)
}

func (router *Router) Patch(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, PATCH, handlers)
}

func (router *Router) Head(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, HEAD, handlers)
}

func (router *Router) Options(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, OPTIONS, handlers)
}

func (router *Router) Any(path string, handlers ...func(http.Handler) http.Handler) {
	router.set(path, "", handlers)
}

func (router *Router) Dir(path string, dir string, handlers ...func(http.Handler) http.Handler) {
	fs := http.FileServer(http.Dir(dir))

	for _, handler := range handlers {
		fs = handler(fs)
	}

	path = "/" + strings.Trim(path, "/") + "/"
	router.mux.Handle("GET "+path, http.StripPrefix(path, fs))
}

func (router *Router) Listen(port string) error {

	return http.ListenAndServe(":"+port, router.mux)
}
