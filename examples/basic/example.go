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
        // add an endpoint to our server
        app.AddRoute("/", hello_world)
        // start a HTTP-only web server
        app.Serve()
}
