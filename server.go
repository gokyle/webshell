package webshell

import (
	"crypto/tls"
        "log"
	"net/http"
	"time"
)

const WEBSHELL_VERSION = "2.0.0"

var (
	ReadTimeout  = 3 * time.Second
	WriteTimeout = 3 * time.Second
)

type WebApp struct {
	name string
	host string
	port string
	key  string
	cert string
	srv  *http.Server
	mux  *http.ServeMux
}

// Retrieve the app's name.
func (app *WebApp) Name() string {
	return app.name
}

// Retrieve the app's host.
func (app *WebApp) Host() string {
	return app.host
}

// Retrieve the app's port.
func (app *WebApp) Port() string {
	return app.port
}

// Retrieve the address (host + port) the app will serve on.
func (app *WebApp) Address() string {
	return serverAddress(app.host, app.port)
}

// IsTLS returns true if the app is configured for TLS.
func (app *WebApp) IsTLS() bool {
	return app.srv.TLSConfig != nil
}

func serverAddress(host, port string) string {
	if host == "" {
		return ":" + port
	} else {
		return host + ":" + port
	}
	return ""
}

func NewApp(name, host, port string) *WebApp {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:           serverAddress(host, port),
		Handler:        mux,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	app := WebApp{name, host, port, "", "", srv, mux}
	return &app
}

func NewTLSApp(name, host, port, key, cert string) *WebApp {
	app := NewApp(name, host, port)
	app.key = key
	app.cert = cert
	app.srv.TLSConfig = new(tls.Config)
	return app
}

// Serve enables the web server; if it is configured for TLS, it will
// listen for and serve TLS requests. If not, it will serve standard
// HTTP requests.
func (app *WebApp) Serve() error {
        log.Println("now listening on ", app.Address())
	if app.IsTLS() {
		return app.srv.ListenAndServeTLS(app.cert, app.key)
	} else {
		return app.srv.ListenAndServe()
	}
	return nil
}
