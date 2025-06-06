package proxy

import (
	"log"
	"net/http/httputil"
	"net/url"
)

// NewSingleHostReverseProxy — обычный HTTP reverse proxy
func NewSingleHostReverseProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid proxy target: %s — %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

// NewWebSocketReverseProxy — WebSocket proxy (на основе обычного)
func NewWebSocketReverseProxy(target string) *httputil.ReverseProxy {
	// Здесь можно будет потом вставить custom Director, Upgrader, Logger и т.п.
	return NewSingleHostReverseProxy(target)
}
