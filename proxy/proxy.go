package main

import (
    "bytes"
    "fmt"
    "os"
    "io"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "net/http/cookiejar"
    "crypto/x509"
    "crypto/tls"
    "net/url"
    "strings"
    "time"
)

func NewProxy(target string) (*Proxy, error) {
    targetUrl, err := url.Parse(target)
    if err != nil {
        return nil, err
    }
   
    // http over tls 
    config    := &tls.Config{
        RootCAs: loadCA("ca.crt"),
    }
    transport := &http.Transport{
        TLSClientConfig: config,
    }
    jar, err  := cookiejar.New(nil)
    if err != nil {
        return nil, err
    }

    // create proxy server
    p := Proxy{}
    p.target = targetUrl
    p.client = &http.Client{
        Transport: transport,
        Jar:       jar, 
    }
    return &p, nil
}

type Proxy struct {
    // url of the server
    target *url.URL

    // proxy server acts as a client to the server
    client *http.Client 
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Do processing of requests here 
    logIP(r)
    
    // example: non curl users are bots
    // do not forward to server
    if !strings.Contains(r.UserAgent(), "curl") {
        io.WriteString(w, "you are a bot\n")
        return
    }

    // get from server
    resp, err := p.client.Get(p.target.String())
    if err != nil {
        log.Println(err)
    }

    // forward response to client
    bytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
    }
    w.Write(bytes)

    // inject a data into the body
    io.WriteString(w, "proxy says hello\n")
}

// change to tick based approach
func (p *Proxy) refresh(delay time.Duration) {
    url := fmt.Sprintf("%s/login", p.target)
    for {
        if delay > 0 {
            time.Sleep(delay)
        }

        // refresh session 
        _, err := p.client.Get(url)
        if err != nil {
            log.Println(err)
        }
    }
}

func (p *Proxy) login() {
    // send password to server for authentication
    url := fmt.Sprintf("%s/login", p.target)
    data := bytes.NewBuffer([]byte(os.Getenv("SERVER_PASS")))
    r, err := p.client.Post(url, "text/html", data)
    if err != nil {
        log.Println(err)
    }
    log.Println(r.Cookies())

    // todo handle bad request
    if r.StatusCode == http.StatusBadRequest {
        return
    }
  
    // refresh session automatically
    cookies := r.Cookies()
    if len(cookies) < 1 || cookies[0].Name != "session_token" {
        log.Println("No session cookie found")
        return
    }
    delay := cookies[0].Expires.Sub(time.Now().Add(2 * time.Second))
    go p.refresh(delay)
}

func logIP(r *http.Request) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil {
        log.Println(net.ParseIP(ip))  
    }
}

func loadCA(caFile string) *x509.CertPool {
    pool := x509.NewCertPool()
    ca, err := ioutil.ReadFile(caFile)
    if  err != nil {
        log.Fatal(err)
    }
    pool.AppendCertsFromPEM(ca)
    return pool
}

func main() {
    proxy, err := NewProxy(os.Getenv("SERVER_URL"))
    if err != nil {
        log.Fatal(err)
    }

    s := &http.Server{
        Addr: ":5000",
        Handler: proxy,
        TLSConfig: &tls.Config{},
    }
    log.Println("Listening :443")

    // authenticate
    proxy.login()
    log.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))
}

