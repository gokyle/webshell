## webshell
### a shell for new Go webapps

### Introduction

`webshell` is a simple framework for quickly getting started with new
webapps in Go that can be configured for TLS or insecure operation.

```go
// example/example.go: very quick example program
package main

import (
        "github.com/gokyle/webshell"
        "log"
        "net/http"
)

func hello_world(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hello, " + r.RemoteAddr))
}

func main() {
        // create the app
        app := webshell.NewApp("example app", "127.0.0.1", "8080")
        // add an endpoint to our server
        app.AddRoute("/", hello_world)
        // start a HTTP-only web server
        log.Fatal(app.Serve())
}
```

### Creating a New WebApp

There are two ways to create a new webapp:

* `NewApp(name, host, port) *WebApp` creates an HTTP app
* `NewTLSApp(name, host, port, keypath, certpath) *WebApp` creates a new TLS
app.

The returned app can be started with `Serve()`; an app can be queried for its
name, host and port with the `Name()`, `Host()`, and `Port()` methods. It can
also be queried to determine if it is a TLS app using the `IsTLS()` method.

### Adding Routes

WebApps have three methods for adding new routes:

* `AddRoute(route, handler)` will add a new route, panicking if the route
couldn't be added.
* `AddConditionalRoute(condition, route, path)` adds the route if condition
is true; it, too, will panic on error.
* `StaticRoute(route, path)` runs a basic file server on the directory, i.e.
for static assets.

### Examples
Contained in the `examples` subdirectory:
* `basic`: bare bones example
* `templates`: templating example
* `static`: demonstrates the use of the static serving functions

Each example should be run from its respective directory, as some use
relative paths in their routes.

### Subpackages
* `webshell/assetcache` provides a simple file cache for static assets
that can help speed up asset delivery.
* `webshell/auth` provides password authentication code; the user need only
supply a function that translates a user ID into a pair of byte slices. It
uses the PBKDF2 key derivation function.

### Under Development
* `webshell/logging` will provide a logging interface for requests.

### License

`webshell` is licensed under an ISC license. The `LICENSE` file contains
the full text of the license.
