// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package proxy contains helper middleware such as request logging.
package proxy

import (
	"log"
	"net/http"

	"DeadEndProxy/config"
)

// logRequests logs each incoming request and the target backend.
func logRequests(next http.Handler, loc config.LocationConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("➡️  %s %s → %s", r.Method, r.URL.Path, loc.ProxyPass)
		next.ServeHTTP(w, r)
	})
}
