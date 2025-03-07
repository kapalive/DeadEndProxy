package security

import "net/http"

func FilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TarpitMiddleware(next http.Handler) http.Handler {
	return next
}

func FakeErrorMiddleware(next http.Handler) http.Handler {
	return next
}
