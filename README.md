# nanoproxy

This is a tiny HTTP forward proxy written in [Go],
for me to gain experience in the Go language.

This proxy accepts all requests and forwards them directly to the
origin server. It performs no caching.

Despite this not being a full proxy implementation, it is blazing fast.
In particular it is significantly faster than Squid and slightly faster than
Apache's mod_proxy. This demonstrates that Go's built-in HTTP library is
of a very high quality and that the Go runtime is quite performant.

Only `xkcd.com` has been really tested with this proxy.
Many other sites don't work with the current implementation.

## Prerequisites

* Go 1.3.3, or a compatible version

## Installation

* Clone this repository.

```
git clone git@github.com:davidfstr/nanoproxy.git
cd nanoproxy
```

* Configure your web browser to route all HTTP traffic through `localhost:8080`.

## Usage

* Start the proxy: `go run nanoproxy.go`

* Open your web browser to `http://xkcd.com` or some other page on that site.

[Go]: https://golang.org

## Notes

* Go's HTTP server implementation is really good. I read it all.
  Only missing feature I desire is the ability to process multiple
  pipelined HTTP requests in parallel.
* Go's HTTP client implementation is easy to use, based on my limited
  experience in this proxy. I have not read its implementation.