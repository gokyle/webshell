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
        // create the asset cache, caching 4 files up to 4MB
        assetcache.MaxItems = 4
        assetcache.MaxSize = 4194304
        ac := assetcache.CreateAssetCache("/assets/", "assets/")
        if err := ac.Start(); err != nil {
                panic("could not start asset cache: " + err.Error())
        }
        app.AddRoute("/assets/", assetcache.AssetHandler(ac))
        app.StaticRoute("/", "static")
        // start a HTTP-only web server
        log.Fatal(app.Serve())
}
