package main

import (
        "fmt"
        "github.com/gokyle/webshell"
        "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
        page := "Check out /test.html and /test2.html ..."
        w.Write([]byte(page))
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
                webshell.Error500(err.Error(), "text/plain", w, r)
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
                webshell.Error500(err.Error(), "text/plain", w, r)
        } else {
                w.Write(out)
        }
}

func tpl_error(w http.ResponseWriter, r *http.Request) {
        var page struct {
                Nonsense string
        }
        out, err := webshell.ServeTemplate("templates/test.html", page)
        if err != nil {
                fmt.Println("[!] error")
                webshell.Error500(err.Error(), "text/plain", w, r)
        } else {
                w.Write(out)
        }
}

func main() {
        // load the requisite environment variables
        webshell.LoadEnv()
        // add an endpoint to our server
        webshell.AddRoute("/", index)
        webshell.AddRoute("/test.html", tpl_test)
        webshell.AddRoute("/test2.html", tpl_test2)
        webshell.AddRoute("/error.html", tpl_error)
        // start a HTTP-only web server
        webshell.Serve(false, nil)
}
