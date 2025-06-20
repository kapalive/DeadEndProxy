// © 2023 Devinsidercode CORP. Licensed under the MIT License.
package security

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Middleware: Sets secure HTTP headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=()")
		next.ServeHTTP(w, r)
	})
}

// Middleware: Bearer token verification (primitive)
func RequireBearerAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(444) // nginx-style custom status
			_, _ = w.Write([]byte("Missing or invalid Bearer token"))
			log.Printf("🚫 Unauthorized request to %s", r.URL.Path)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// TarpitMiddleware: Simulating Latency for Bad Requests
func TarpitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// В будущем можно проверять IP, User-Agent и т.п.
		time.Sleep(1 * time.Second)
		next.ServeHTTP(w, r)
	})
}

// FakeErrorMiddleware: Random Error Substitution (Obfuscation Attacks)
func FakeErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Тут можно случайно вернуть 403, 404, 500 — обманка
		// или подменять заголовки на несуществующие
		w.Header().Set("Server", "Devinsider Proxy/1.0.0 (Ubuntu)")
		next.ServeHTTP(w, r)
	})
}

// FilterMiddleware: Filter example (blank)
func FilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// В будущем сюда пойдут geo-block, referer-check, IP denylist и т.п.
		next.ServeHTTP(w, r)
	})
}

// ApplySecurityChain — wraps the handler with all the middleware at once
func ApplySecurityChain(h http.Handler, withBearer bool) http.Handler {
	chain := h
	chain = FilterMiddleware(chain)
	chain = FakeErrorMiddleware(chain)
	chain = TarpitMiddleware(chain)
	chain = SecurityHeadersMiddleware(chain)
	if withBearer {
		chain = RequireBearerAuthorization(chain)
	}
	return chain
}
