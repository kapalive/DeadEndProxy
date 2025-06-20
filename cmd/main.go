package main

import (
	"DeadEndProxy/config"
	"DeadEndProxy/internal/router"
	"flag"
	"fmt"
	"os"

	"DeadEndProxy/internal/proxy"

	"github.com/redis/go-redis/v9"
)

func main() {
	portHTTP := flag.String("port-http", "80", "HTTP port")
	portHTTPS := flag.String("port-proxy", "443", "HTTPS port")
	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println("Usage of DeadEndProxy:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	if *help {
		fmt.Println("Usage of DeadEndProxy:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	config.MustLoadInitial(*configPath)
	config.WatchAndReload(*configPath)

	redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	resolver := router.NewResolver(redisClient)

	proxy.StartWithOverride(&proxy.ConfigOverride{
		HTTPPort:  *portHTTP,
		HTTPSPort: *portHTTPS,
	}, resolver)

}
