package webshell

import (
	"log"
	"net/http"
	"path/filepath"
)

// Simple name for a function capable of handling a route.
type RouteHandler func(w http.ResponseWriter, r *http.Request)

// Simple name for a function that handles basic errors.
type ErrorRoute func(string, string, http.ResponseWriter, *http.Request)
type TemplateErrorRoute func(interface{}, http.ResponseWriter, *http.Request)

// Generic error handlers. They take a message and content-type as a string,
// as well as the HTTP response writer and request, and respond with the
// named error. The http.StatusText function may be used to return the
// text for the error code. Note that you have to explicitly call these; the
// design of Go's http server means it will respond with its own 404 handler
// if a route is not found or if a static file server cannot find the file.
var (
	Error400 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error401 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error403 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error404 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error405 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error429 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error500 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error501 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error502 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
	Error503 func(msg, ctype string, w http.ResponseWriter, r *http.Request)
)

func (app *WebApp) AddRoute(path string, handler RouteHandler) {
	app.mux.HandleFunc(path, handler)
	log.Printf("[+] route %s added\n", path)
}

func (app *WebApp) AddConditionalRoute(condition bool, path string, handler RouteHandler) {
	if condition {
		app.AddRoute(path, handler)
	}
}

func (app *WebApp) StaticRoute(route string, path string) {
	var err error
	if len(route) == 0 {
		panic("Invalid route: " + route + " -> " + path)
	}
	if len(path) == 0 {
		panic("Invalid path:" + route + " -> " + path)
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			panic(err)
		}
	}
	app.mux.Handle(route, http.StripPrefix(route, http.FileServer(http.Dir(path))))
	log.Printf("static route %s -> %s added\n", route, path)
}

// GenerateErrorHandler returns a RouteHandler function
func GenerateErrorHandler(status int) ErrorRoute {
	return func(msg, ctype string, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Add("content-type", ctype)
		w.Write([]byte(msg))
	}
}

// GenerateTemplateErrorHandler returns a function serving a templated error
func GenerateTemplateErrorHandler(status int, filename string) (hdlr TemplateErrorRoute, err error) {
	tpl, err := CompileTemplate(filename)
	if err != nil {
		return
	}
	hdlr = func(in interface{}, w http.ResponseWriter, r *http.Request) {
		msg, err := ServeTemplate(tpl, in)
		if err != nil {
			log.Printf("error serving template %d %s: %s\n",
				status, filename, err.Error())
			return
		}
		w.WriteHeader(status)
		w.Write(msg)
	}
	return
}

func init() {
	Error400 = GenerateErrorHandler(http.StatusBadRequest)
	Error401 = GenerateErrorHandler(http.StatusUnauthorized)
	Error403 = GenerateErrorHandler(http.StatusForbidden)
	Error404 = GenerateErrorHandler(http.StatusNotFound)
	Error429 = GenerateErrorHandler(429)
	Error500 = GenerateErrorHandler(http.StatusInternalServerError)
	Error501 = GenerateErrorHandler(http.StatusNotImplemented)
	Error502 = GenerateErrorHandler(http.StatusBadGateway)
	Error503 = GenerateErrorHandler(http.StatusServiceUnavailable)

}

func ContentResponder(r *http.Request) string {
	accept := r.Header["Accept"][0]
	if accept == "" {
		return "text/plain"
	} else if accept == "*/*" {
		return "text/plain"
	}
	return accept
}
