package proxy

import (
	"DeadEndProxy/internal/router"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"DeadEndProxy/config"
	"DeadEndProxy/internal/security"
)

// Main Start
func Start(cfg *config.Config, resolver *router.Resolver) {
	startHTTPRedirect(cfg)
	startHTTPSProxy(cfg, resolver)
}

// Start with CLI
func StartWithOverride(override *ConfigOverride, resolver *router.Resolver) {
	cfg := config.MustLoadConfig("config.yaml")
	override.Apply(cfg)
	Start(cfg, resolver)

	if override != nil {
		if override.HTTPPort != "" {
			log.Printf("Overriding HTTP port: %d -> %s", cfg.Server.HTTPPort, override.HTTPPort)
			override.Apply(cfg)
		}
		if override.HTTPSPort != "" {
			log.Printf("Overriding HTTPS port: %d -> %s", cfg.Server.HTTPSPort, override.HTTPSPort)
			override.Apply(cfg)
		}
	}

	Start(cfg, resolver)
}

// 🔐 HTTPS с TLS + SNI + ReverseProxy
func startHTTPSProxy(cfg *config.Config, resolver *router.Resolver) {
	mainCert, err := tls.LoadX509KeyPair(cfg.Server.SSLCertMain, cfg.Server.SSLKeyMain)
	if err != nil {
		log.Fatal("TLS Main cert load failed:", err)
	}
	secondCert, err := tls.LoadX509KeyPair(cfg.Server.SSLCertSecond, cfg.Server.SSLKeySecond)
	if err != nil {
		log.Fatal("TLS Second cert load failed:", err)
	}

	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
		Certificates: []tls.Certificate{mainCert, secondCert},
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if hello == nil {
				return &mainCert, nil
			}
			host := strings.ToLower(hello.ServerName)
			if strings.Contains(host, strings.ToLower(cfg.Server.DomainSecond)) {
				return &secondCert, nil
			}
			return &mainCert, nil
		},
	}

	addrHTTPS := ":" + strconv.Itoa(cfg.Server.HTTPSPort)
	server := &http.Server{
		Addr:      addrHTTPS,
		TLSConfig: tlsCfg,
		Handler:   buildRootHandler(cfg, resolver),
	}

	log.Printf("Starting HTTPS proxy on %s ...", addrHTTPS)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal("HTTPS server failed:", err)
	}
}

// (root)
func buildRootHandler(cfg *config.Config, resolver *router.Resolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domain := strings.ToLower(r.Host)
		handler := buildDomainHandler(cfg, domain)

		if handler == nil {
			// 🚨 Подключаем динамическую маршрутизацию
			handler = createDynamicHandler(resolver)
		}

		handler.ServeHTTP(w, r)
	})
}

func buildDomainHandler(cfg *config.Config, domain string) http.Handler {
	mux := http.NewServeMux()

	// Static для главного домена
	if strings.EqualFold(domain, cfg.Server.DomainMain) {
		fs := http.FileServer(http.Dir(cfg.Server.Webroot))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := cfg.Server.Webroot + r.URL.Path
			if _, err := os.Stat(path); err != nil {
				http.ServeFile(w, r, cfg.Server.Webroot+"/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		})
	}

	// Прокси по локациям
	for _, loc := range cfg.Server.Locations {
		if strings.EqualFold(loc.Domain, domain) {
			handler := createProxyHandler(loc)
			mux.Handle(loc.Path, handler)
		}
	}

	return mux
}

func createProxyHandler(loc config.LocationConfig) http.Handler {
	var handler http.Handler

	if loc.IsWebSocket {
		handler = NewWebSocketReverseProxy(loc.ProxyPass)
	} else {
		handler = NewSingleHostReverseProxy(loc.ProxyPass)
	}

	// Apply security chain
	handler = security.ApplySecurityChain(handler, loc.RequireBearer)

	// Apply CORS if needed (после безопасности, до ServeHTTP)
	if loc.Cors {
		handler = HandleCORS(handler)
	}

	return handler
}
