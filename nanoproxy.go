package main

import (
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
    passthruRequestHeaderKeys := [...]string{
        "Accept",
        "Accept-Encoding",
        "Accept-Language",
        "Cache-Control",
        "Cookie",
        "Referer",
        "User-Agent",
    }
    
    passthruResponseHeaderKeys := [...]string{
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
    
    handler := http.DefaultServeMux
    
    s := &http.Server{
        Addr:           ":8080",
        Handler:        handler,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    
    handleFunc := func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("request from client: %+v\n", r)
        
        hh := http.Header{}
        for _, hk := range passthruRequestHeaderKeys {
            if hv, ok := r.Header[hk]; ok {
                hh[hk] = hv
            }
        }
        
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
        
        fmt.Printf("response from origin server: %+v\n", resp)
        
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
    
    // TODO: There's no way to register a handler that listens to ALL
    //       incoming requests at the moment. Good idea to ask the
    //       Go community for a workaround.
    handler.HandleFunc("xkcd.com/", handleFunc)
    handler.HandleFunc("imgs.xkcd.com/", handleFunc)
    handler.HandleFunc("www.google.com/", handleFunc)
    handler.HandleFunc("c.xkcd.com/", handleFunc)
    
    s.ListenAndServe()
}