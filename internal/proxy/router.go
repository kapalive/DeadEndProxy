package proxy

import (
	"log"
	"net/http"

	"DeadEndProxy/config"
)

// logRequests — middleware для логов
func logRequests(next http.Handler, loc config.LocationConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("➡️  %s %s → %s", r.Method, r.URL.Path, loc.ProxyPass)
		next.ServeHTTP(w, r)
	})
}
