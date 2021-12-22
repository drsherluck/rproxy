package main

import (
    "os"
    "log"
    "net/http"
    "fmt"
    "io"
    "time"
    "github.com/google/uuid"
    c "github.com/patrickmn/go-cache"
)

const (
    // how long a session is valid
    sessionExpireTime = 60 * time.Second
    // the key of the session cookie
    sKey = "session_token"
)

var (
    // password to needed to access this server
    password string
    // the session cache
    cache    *c.Cache
)

type handler struct {
    counter int
}

// default handler
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    _, ok := checkSessionValidity(w, r)
    if ok {
        h.counter++
        fmt.Fprintf(w, "server handled %d requests\n", h.counter)
    }
}

// handler for client autentication
func authenticate(w http.ResponseWriter, r *http.Request) {
    // request authorization
    if r.Method == http.MethodPost {
        pass, err := io.ReadAll(r.Body)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
    
        // check that passwords match and create session
        if password != string(pass) {
            log.Println("Unsuccessful login attempt")
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        createSession(w, r)
    } else { 
        // refresh session
        session, ok := checkSessionValidity(w, r)
        if ok {
            cache.Delete(session)
            createSession(w, r)
        }
    }
}

// checks if a session exist and if it is valid
// writes appropiate status codes to the response header on bad request and 
// unauthorized access
func checkSessionValidity(w http.ResponseWriter, r *http.Request) (session string, ok bool) {
    // retreive cookie
    cookie, err := r.Cookie(sKey)
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
        } else {
            w.WriteHeader(http.StatusBadRequest)
        }
        return "", false
    }

    // check if session is valid
    session = cookie.Value
    _, found := cache.Get(session)
    if !found {
        w.WriteHeader(http.StatusUnauthorized)
    }
    return session, found
}

// create a new session id and set the response cookie
func createSession(w http.ResponseWriter, r *http.Request) {
    // create new session id
    session := uuid.New().String()
    cache.Set(session, 0, sessionExpireTime)

    // set the session cookie
    http.SetCookie(w, &http.Cookie{
        Name: sKey,
        Value: session,
        Expires: time.Now().Add(sessionExpireTime),
        //Secure: true,
        //HttpOnly: true,
    })
}

func main() {
    // load password
    password = os.Getenv("SERVER_PASS")
    // setup cache
    cache = c.New(c.NoExpiration, time.Hour)

    // creat server and route handlers
    mux := http.NewServeMux()
    mux.Handle("/", &handler{0})
    mux.HandleFunc("/login", authenticate)

    s := &http.Server {
        Addr: ":80",
        Handler: mux,
    }
    log.Println("Listening :80")
    log.Fatal(s.ListenAndServe());
}
