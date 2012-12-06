package main

import (
        "fmt"
	"github.com/gokyle/webshell"
	"github.com/gokyle/webshell/auth"
        "net/http"
)

var (
	app   *webshell.WebApp
	store *auth.SessionStore
        index = webshell.MustCompileTemplate("templates/index.html")
)

type Page struct {
        Authenticated bool
        Status int
        Cookie *http.Cookie
}

func init() {
        hash, salt := auth.HashPass("hello, world")
        SetUserPass("joe", hash, salt)
}

func main() {
	auth.LookupCredentials = LookupUser
	app = webshell.NewApp("authentication tests", "", "9000")
	store = auth.CreateSessionStore(`s_nm`, false, nil)
	app.AddRoute("/", testAuth)
        app.StaticRoute("/assets/", "assets/")
	app.Serve()
}

func testAuth(w http.ResponseWriter, r *http.Request) {
        var page Page
        if store.CheckSession(r) {
                page.Authenticated = true
                page.Status = http.StatusOK
        } else {
                page.Status = http.StatusUnauthorized
        }
        if r.Method == "POST" {
                processForm(w, r)
                return
        }

        serveIndex(page, w, r)
}

func serveIndex(page Page, w http.ResponseWriter, r *http.Request) {
        cookie := page.Cookie
        page.Cookie = nil
        out, err := webshell.ServeTemplate(index, page)
        if err != nil {
                webshell.Error500(err.Error(), "text/plain", w, r)
                return
        }
        if cookie != nil {
                fmt.Println("[+] setting cookie")
                fmt.Printf("\t[*] cookie: %+v\n", cookie)
                http.SetCookie(w, cookie)
        }
        w.Write(out)
}

func processForm(w http.ResponseWriter, r *http.Request) {
        var page Page
        page.Status = http.StatusUnauthorized
        err := r.ParseForm()
        if err != nil {
                webshell.Error400(err.Error(), "text/plain", w, r)
                return
        }

        user := r.Form.Get("user")
        pass := r.Form.Get("pass")
        logout := r.Form.Get("logout")
        if user == "" || pass == "" {
                if logout == "true" {
                        store.DestroySession(r)
                        page.Status = http.StatusUnauthorized
                }
                serveIndex(page, w, r)
                return
        }

        cookie := store.AuthSession(user, pass, r, nil)
        if cookie == nil {
                serveIndex(Page{}, w, r)
                return
        }

        page.Cookie = cookie
        page.Authenticated = true
        page.Status = http.StatusOK
        serveIndex(page, w, r)
}
