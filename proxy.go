package main

import (
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func NewProxy(target string) (*Proxy, error) {
    targetUrl, err := url.Parse(target)
    if err == nil {
        p := Proxy{}
        p.proxy   = httputil.NewSingleHostReverseProxy(targetUrl)
        return &p, nil
    }
    return nil, err
}

type Proxy struct {
    proxy *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Do processing of requests here 
    logIP(r)
    p.proxy.ServeHTTP(w, r) 
}

func logIP(r *http.Request) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil {
        log.Println(net.ParseIP(ip))  
    }
}

func main() { 
    // Create proxy for the server 
    proxy, err := NewProxy("http://localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    
    s := &http.Server {
        Addr: ":5000", // proxy url
        Handler: proxy,
    }
    log.Println("Listening http://localhost:5000")
    log.Fatal(s.ListenAndServe())
}

