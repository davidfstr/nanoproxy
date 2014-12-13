package main

import (
    "fmt"
    "io"
    "net/http"
    "time"
)

var verbose = false;

var passthruRequestHeaderKeys = [...]string{
    "Accept",
    "Accept-Encoding",
    "Accept-Language",
    "Cache-Control",
    "Cookie",
    "Referer",
    "User-Agent",
}

var passthruResponseHeaderKeys = [...]string{
    "Content-Encoding",
    "Content-Language",
    "Content-Type",
    "Cache-Control", // TODO: Is this valid in a response?
    "Date",
    "Etag",
    "Expires",
    "Last-Modified",
    "Location",
    "Server",
    "Vary",
}

func main() {
    handler := http.DefaultServeMux
    
    handler.HandleFunc("/", handleFunc)
    
    s := &http.Server{
        Addr:           ":8080",
        Handler:        handler,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    
    s.ListenAndServe()
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("--> %v %v\n", r.Method, r.URL)
    
    // Construct filtered header to send to origin server
    hh := http.Header{}
    for _, hk := range passthruRequestHeaderKeys {
        if hv, ok := r.Header[hk]; ok {
            hh[hk] = hv
        }
    }
    
    // Construct request to send to origin server
    rr := http.Request{
        Method: r.Method,
        URL: r.URL,
        Header: hh,
        Body: r.Body,
        // TODO: Is this correct for a 0 value?
        //       Perhaps a 0 may need to be reinterpreted as -1?
        ContentLength: r.ContentLength,
        Close: r.Close,
    }
    
    // Forward request to origin server
    resp, err := http.DefaultTransport.RoundTrip(&rr)
    if err != nil {
        // TODO: Passthru more error information
        http.Error(w, "Could not reach origin server", 500)
        return
    }
    defer resp.Body.Close()
    
    if (verbose) {
        fmt.Printf("<-- %v %+v\n", resp.Status, resp.Header)
    } else {
        fmt.Printf("<-- %v\n", resp.Status)
    }
    
    // Transfer filtered header from origin server -> client
    respH := w.Header()
    for _, hk := range passthruResponseHeaderKeys {
        if hv, ok := resp.Header[hk]; ok {
            respH[hk] = hv
        }
    }
    w.WriteHeader(resp.StatusCode)
    
    // Transfer response from origin server -> client
    if resp.ContentLength > 0 {
        // (Ignore I/O errors, since there's nothing we can do)
        io.CopyN(w, resp.Body, resp.ContentLength)
    } else if (resp.Close) { // TODO: Is this condition right?
        // Copy until EOF or some other error occurs
        for {
            if _, err := io.Copy(w, resp.Body); err != nil {
                break
            }
        }
    }
}