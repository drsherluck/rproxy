package main

import (
    "io"
    "log"
    "net/http"
    "fmt"
)

type handler struct {
    counter int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.counter++
    io.WriteString(w, fmt.Sprintf("server handled %d requests", h.counter))
}

func main() {
    s := &http.Server {
        Addr: ":8080",
        Handler: &handler{0},
    }
    log.Println("Listening http://localhost:8080")
    log.Fatal(s.ListenAndServe())    
}
