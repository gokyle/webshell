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
        webshell.LoadEnv()
        // add an endpoint to our server
        webshell.AddRoute("/", hello_world)
        // start a HTTP-only web server
        webshell.Serve(false, nil)
}
