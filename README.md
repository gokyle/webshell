## webshell
### a shell for new Go webapps

### Introduction

`webshell` is a simple framework for quickly getting started with new
webapps in Go. It loads all of its configuration from environment
variables, and can be configured for TLS or insecure operation.

```go
// example/example.go: very quick example program
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
        webshell.AddRoute("/hello", hello_world)
        // start a HTTP-only web server
        webshell.Serve(false, nil)
}
```

### Configuration

`webshell` by default pulls its configuration in from the environment.
`LoadEnv()` will set the relevant variables. However, this can be
bypassed to use your own configuration method. The relevant variables
are (note that the environment variables have the same name):

* `SERVER_ADDR` contains the address the server should listen on.
* `SERVER_PORT` contains the port the server should listen on.
* `SSL_KEY` contains the path to the SSL private key.
* `SSL_CERT` contains the path to the SSL certficate.

For example, to load the server address from a function called
`LoadAddressFromDB`:

```go
webshell.SERVER_ADDR = LoadAddressFromDB()
```

An example shell script that can be sourced to sane defaults for the
server may be found in `examples/env.sh`.

### License

`webshell` is licensed under an ISC license. The `LICENSE` file contains
the full text of the license.

