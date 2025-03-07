package main

import (
	"DeadEndProxy/config"
	"DeadEndProxy/internal/proxy"
	"DeadEndProxy/internal/security"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Создаём сервер
	mux := http.NewServeMux()

	// Добавляем маршруты из конфига
	for _, route := range cfg.Proxy.Routes {
		mux.HandleFunc(route.Path, proxy.DynamicRouter(route.Backend))
	}

	// Включаем защитные механизмы
	securedMux := security.FilterMiddleware(mux)
	securedMux = security.TarpitMiddleware(securedMux)
	securedMux = security.FakeErrorMiddleware(securedMux)

	log.Println("💀 DeadEndProxy running on port 8080...")

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), securedMux)
	if err != nil {
		log.Fatalf("Server startup error: %v", err)
	}

	fmt.Println(`
██████╗ ███████╗ █████╗ ██████╗ ███████╗███╗   ██╗██████╗ ███████╗███╗   ██╗
██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝████╗  ██║██╔══██╗██╔════╝████╗  ██║
██║  ██║█████╗  ███████║██║  ██║███████╗██╔██╗ ██║██║  ██║█████╗  ██╔██╗ ██║
██║  ██║██╔══╝  ██╔══██║██║  ██║╚════██║██║╚██╗██║██║  ██║██╔══╝  ██║╚██╗██║
██████╔╝███████╗██║  ██║██████╔╝███████║██║ ╚████║██████╔╝███████╗██║ ╚████║
╚═════╝ ╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═══╝`)

	http.ListenAndServe(":8080", securedMux)
}
