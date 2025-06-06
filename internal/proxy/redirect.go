package proxy

import (
	"log"
	"net/http"
	"strconv"

	"DeadEndProxy/config"
)

// startHTTPRedirect запускает редирект HTTP → HTTPS
func startHTTPRedirect() {
	go func() {
		cfg := config.GetConfig() // ✅ теперь он сам знает, где взять актуальный конфиг
		addr := ":" + strconv.Itoa(cfg.Server.HTTPPort)
		log.Printf("🌐 HTTP: %s (redirect → HTTPS)", addr)

		redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			log.Printf("➡️  Redirect: %s → %s", r.URL.String(), target)
			http.Redirect(w, r, target, http.StatusPermanentRedirect)
		})

		if err := http.ListenAndServe(addr, redirectHandler); err != nil {
			log.Fatalf("❌ HTTP Redirect failed: %v", err)
		}
	}()
}
