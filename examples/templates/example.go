package main

import (
        "fmt"
        "github.com/gokyle/webshell"
        "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
        page := `<!doctype html>
<html>
  <head>
    <title>webshell template example</title>
  </head>

  <body>
  <h1>webshell template example</h1>
  <p>Take a look at the following pages:</p>
  <ul>
    <li><a href="/test.html">first test page</a></li>
    <li><a href="/test2.html">second test page</a></li>
    <li><a href="/error.html">template error page</a></li>
   </ul>
  </body>
</html>
`
        w.Write([]byte(page))
}

func tpl_test(w http.ResponseWriter, r *http.Request) {
        var page struct {
                Title string
                Paragraph string
        }
        page.Title = "test page"
        page.Paragraph = "rw nw prt m hrw"
        out, err := webshell.ServeTemplateFile("templates/test.html", page)
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
        out, err := webshell.ServeTemplateFile("templates/test.html", page)
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
        out, err := webshell.ServeTemplateFile("templates/test.html", page)
        if err != nil {
                fmt.Println("[!] error")
                webshell.Error500(err.Error(), "text/plain", w, r)
        } else {
                w.Write(out)
        }
}

func main() {
        // load the requisite environment variables
        app := webshell.NewApp("webshell basic example", "127.0.0.1", "8080")
        // add an endpoint to our server
        app.AddRoute("/", index)
        app.AddRoute("/test.html", tpl_test)
        app.AddRoute("/test2.html", tpl_test2)
        app.AddRoute("/error.html", tpl_error)
        // start a HTTP-only web server
        app.Serve()
}
