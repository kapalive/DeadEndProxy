// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package main contains the CLI entry point for DeadEndProxy.
// It parses command line flags, loads the configuration file and
// starts the reverse proxy server.
package main

import (
	"DeadEndProxy/config"
	"DeadEndProxy/internal/router"
	"flag"
	"fmt"
	"net/http"
	"os"

	"DeadEndProxy/internal/proxy"

	"github.com/redis/go-redis/v9"
)

// main initializes the configuration and starts the proxy based on
// command line overrides.
func main() {
	portHTTP := flag.String("port-http", "80", "HTTP port")
	portHTTPS := flag.String("port-proxy", "443", "HTTPS port")
	configPath := flag.String("config", "/etc/deadendproxy/config.yaml", "Path to config file")
	help := flag.Bool("h", false, "Show help")

	flag.Parse() // ✅ должен вызываться ТОЛЬКО один раз и после ВСЕХ flag.String

	if *help {
		fmt.Println("Usage of DeadEndProxy:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Загружаем конфиг
	config.MustLoadInitial(*configPath)
	config.WatchAndReload(*configPath)

	// Redis
	redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	resolver := router.NewResolver(redisClient)

	// Статика
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("webroot/static"))))

	// Запуск прокси
	proxy.StartWithOverride(&proxy.ConfigOverride{
		HTTPPort:  *portHTTP,
		HTTPSPort: *portHTTPS,
	}, resolver)
}
