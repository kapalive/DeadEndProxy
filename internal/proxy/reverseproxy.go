// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package proxy wraps the Go standard library reverse proxy to
// provide custom error handling and WebSocket support.
package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewSingleHostReverseProxy creates a basic HTTP reverse proxy
// with custom timeout error handling.
func NewSingleHostReverseProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid proxy target: %s — %v", target, err)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		writeErrorPage(w, http.StatusGatewayTimeout)
	}
	return rp
}

// NewWebSocketReverseProxy creates a reverse proxy suitable for
// WebSocket connections.
func NewWebSocketReverseProxy(target string) *httputil.ReverseProxy {
	// Here you can then insert a custom Director, Upgrader, Logger, etc.
	return NewSingleHostReverseProxy(target)
}
