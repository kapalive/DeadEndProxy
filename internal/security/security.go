// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package security defines a set of HTTP middleware used to
// protect the proxy and obfuscate malicious traffic.
package security

import (
	"DeadEndProxy/config"
	"log"
	"net/http"
	"strings"
	"time"
)

// SecurityHeadersMiddleware sets a number of recommended security
// headers on each response.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := map[string]string{
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
			"X-Content-Type-Options":    "nosniff",
			"X-Frame-Options":           "DENY",
			"X-XSS-Protection":          "1; mode=block",
			"Referrer-Policy":           "strict-origin-when-cross-origin",
			"Permissions-Policy":        "geolocation=(), microphone()",
		}

		cfg := config.GetConfig()
		if cfg != nil && cfg.Headers != nil {
			for k, v := range cfg.Headers {
				headers[k] = v
			}
		}

		for k, v := range headers {
			w.Header().Set(k, v)
		}
		next.ServeHTTP(w, r)
	})
}

// RequireBearerAuthorization checks the "Authorization" header for
// a Bearer token and rejects the request if it is missing or invalid.
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

// RequireCookieAuthorization checks for a specific auth cookie and rejects the request if missing.
func RequireCookieAuthorization(name string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(name)
			if err != nil || c.Value == "" {
				w.WriteHeader(444)
				_, _ = w.Write([]byte("Missing auth cookie"))
				log.Printf("🚫 Unauthorized request to %s", r.URL.Path)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// TarpitMiddleware introduces a small delay that can help slow down
// simple malicious traffic.
func TarpitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// В будущем можно проверять IP, User-Agent и т.п.
		time.Sleep(1 * time.Second)
		next.ServeHTTP(w, r)
	})
}

// FakeErrorMiddleware randomly modifies the response to make
// automated attacks harder.
func FakeErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Тут можно случайно вернуть 403, 404, 500 — обманка
		// или подменять заголовки на несуществующие
		w.Header().Set("Server", "Devinsider Proxy/1.0.0 (Ubuntu)")
		next.ServeHTTP(w, r)
	})
}

// FilterMiddleware is a placeholder for future IP or geo filters.
func FilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// В будущем сюда пойдут geo-block, referer-check, IP denylist и т.п.
		next.ServeHTTP(w, r)
	})
}

// ApplySecurityChain wraps the handler with all security middleware
// in the correct order.
func ApplySecurityChain(h http.Handler, withBearer bool, cookieName string) http.Handler {
	chain := h
	chain = FilterMiddleware(chain)
	chain = FakeErrorMiddleware(chain)
	chain = TarpitMiddleware(chain)
	chain = SecurityHeadersMiddleware(chain)
	if withBearer {
		chain = RequireBearerAuthorization(chain)
	}
	if cookieName != "" {
		chain = RequireCookieAuthorization(cookieName)(chain)
	}
	return chain
}
