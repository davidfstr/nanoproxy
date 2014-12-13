# nanoproxy

This is a tiny HTTP forward proxy written in [Go],
for me to gain experience in the Go language.

Despite this not being a full proxy implementation, it is blazing fast.
In particular it is significantly faster than Squid and slightly faster than
Apache's mod_proxy. This demonstrates that Go's built-in HTTP library is
of a very high quality and that the Go runtime is quite performant.

At the time of writing this proxy is hardcoded to only forward traffic
to `xkcd.com` and subdomains due to API limitations in registering HTTP
handlers in Go's HTTP server library.

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