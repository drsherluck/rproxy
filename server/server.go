package main

import (
    "log"
    "net/http"
    "crypto/tls"
    "fmt"
)

type handler struct {
    counter int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.counter++
    fmt.Fprintf(w, "server handled %d requests\n", h.counter)
}

func main() {
    s := &http.Server {
        Addr: ":8080",
        Handler: &handler{0},
        TLSConfig: &tls.Config{},
    }
    log.Println("Listening :443")
    log.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))    
}
