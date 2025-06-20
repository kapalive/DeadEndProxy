// © 2023 Devinsidercode CORP. Licensed under the MIT License.
package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewSingleHostReverseProxy — обычный HTTP reverse proxy
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

// NewWebSocketReverseProxy — WebSocket proxy (на основе обычного)
func NewWebSocketReverseProxy(target string) *httputil.ReverseProxy {
	// Здесь можно будет потом вставить custom Director, Upgrader, Logger и т.п.
	return NewSingleHostReverseProxy(target)
}
