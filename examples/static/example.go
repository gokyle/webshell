package main

import (
        "github.com/gokyle/webshell"
        "net/http"
)

func hello_world(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hello, " + r.RemoteAddr))
}

func main() {
        // load the requisite environment variables
        app := webshell.NewApp("webshell basic example", "127.0.0.1", "8080")
        // set up our static routes
        app.StaticRoute("/assets/css/", "assets/css/")
        app.StaticRoute("/", "static")
        // start a HTTP-only web server
        app.Serve()
}
