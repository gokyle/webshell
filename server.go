package webshell

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const WEBSHELL_VERSION = "1.0.0"

// Server configuration options.
var (
	SSL_KEY     string
	SSL_CERT    string
	SERVER_ADDR string
	SERVER_PORT string
)

// LoadEnv pulls in configuration data from the environment.
func LoadEnv() {
	SSL_KEY = os.Getenv("SSL_KEY")
	SSL_CERT = os.Getenv("SSL_CERT")
	SERVER_ADDR = os.Getenv("SERVER_ADDR")
	SERVER_PORT = os.Getenv("SERVER_PORT")
}

func main() {
	log.Println("starting server")
	Serve(false, nil)
}

/*
   Server handles the HTTP server set up (with option TLS setup) and
   and serving. It is a blocking function, so any additional functions that
   need to be concurrently run should be fired off via goroutines prior to
   calling Server. If tlsCfg is not nil, the server will use it as the
   TLS configuration for the server. If it is nil, the server will create
   a default configuration. This argument has no effect if the server
   will not be serving TLS requests.
*/
func Serve(doTLS bool, tlsCfg *tls.Config) {
        initDefaultErrors()
	var serverAddress string
	if SERVER_PORT != "" {
		serverAddress = fmt.Sprintf("%s:%s", SERVER_ADDR, SERVER_PORT)
	} else {
		serverAddress = SERVER_ADDR
	}
	log.Println("server address:", serverAddress)
	srv := &http.Server{
		Addr:           serverAddress,
		Handler:        RouterMux,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if doTLS {
		if tlsCfg == nil {
			tlsCfg = new(tls.Config)
		}
		srv.TLSConfig = tlsCfg
		log.Println("listening for incoming TLS connections")
		log.Printf("using credentials:\n\tkey: %s\n\tcert: %s\n",
			SSL_KEY, SSL_CERT)
		log.Fatalf("error in ListenAndServeTLS:\n\t%+v",
			srv.ListenAndServeTLS(SSL_CERT, SSL_KEY))
	} else {
		log.Println("listening for incoming HTTP connections")
		log.Fatalf("error in ListenAndServe:\n\t%+v\n",
			srv.ListenAndServe())
	}
}
