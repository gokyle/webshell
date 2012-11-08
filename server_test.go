package webshell

import (
        "testing"
)

func TestServer(t *testing.T) {
        fmt.Println("This test will run until interrupted with ^c or until")
        fmt.Println("the Go test runner decides the test has run for too")
        fmt.Println("long.")
        LoadEnv()
        Serve(false, nil)
        t.FailNow()
}
