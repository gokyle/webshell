package main

import (
        "github.com/gokyle/webshell/assetcache"
        "github.com/gokyle/webshell"
        "log"
        "net/http"
)

func hello_world(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hello, " + r.RemoteAddr))
}

func main() {
        // load the requisite environment variables
        app := webshell.NewApp("basic asset cache", "", "8080")
        // create the asset cache
        assetcache.BackgroundAttachAssetCache(app, "/assets/",
                "assets/")
        app.StaticRoute("/", "static")
        // start a HTTP-only web server
        log.Fatal(app.Serve())
}
