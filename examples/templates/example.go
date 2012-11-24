package main

import (
        "fmt"
        "github.com/gokyle/webshell"
        "net/http"
)

func hello_world(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hello, " + r.RemoteAddr))
}

func tpl_test(w http.ResponseWriter, r *http.Request) {
        var page struct {
                Title string
                Paragraph string
        }
        page.Title = "test page"
        page.Paragraph = "rw nw prt m hrw"
        out, err := webshell.ServeTemplate("templates/test.html", page)
        if err != nil {
                fmt.Println("[!] error: ", err.Error())
        } else {
                w.Write(out)
        }
}

func tpl_test2(w http.ResponseWriter, r *http.Request) {
        var page struct {
                Title string
                Paragraph string
        }
        page.Title = "another test page"
        page.Paragraph = "Sæmundar Edda"
        out, err := webshell.ServeTemplate("templates/test.html", page)
        if err != nil {
                fmt.Println("[!] error: ", err.Error())
        } else {
                w.Write(out)
        }
}

func main() {
        // load the requisite environment variables
        webshell.LoadEnv()
        // add an endpoint to our server
        webshell.AddRoute("/hello", hello_world)
        webshell.AddRoute("/test.html", tpl_test)
        webshell.AddRoute("/test2.html", tpl_test2)
        // start a HTTP-only web server
        webshell.Serve(false, nil)
}
