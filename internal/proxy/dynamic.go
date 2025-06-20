// © 2023 Devinsidercode CORP. Licensed under the MIT License.
package proxy

import (
	"net/http"
	"strings"

	"DeadEndProxy/internal/router"
)

// createDynamicHandler проксирует по DNS-настроенному кастомному домену
func createDynamicHandler(resolver *router.Resolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.ToLower(r.Host)
		ctx := r.Context()

		entry, err := resolver.ResolveDomain(ctx, host)
		if err != nil || entry == nil {
			writeErrorPage(w, http.StatusBadGateway)
			return
		}

		proxy := NewSingleHostReverseProxy(entry.Target)
		proxy.ServeHTTP(w, r)
	})
}
