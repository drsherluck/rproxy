package main

import (
    "io"
    "log"
    "net/http"
)

func main() {
    resp, err := http.Get("http://localhost:5000")
    if err != nil {
        log.Fatal(err)
    }
    b, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("%s\n", b)
}



