// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package router contains logic for resolving custom domains
// using Redis, DNS TXT records and a fallback core API.
package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RouteEntry struct {
	Target   string `json:"target"`   // Proxy target
	Username string `json:"username"` // For audit
}

type Resolver struct {
	Redis *redis.Client
}

// NewResolver creates a new Resolver instance using the
// provided Redis client.
func NewResolver(redisClient *redis.Client) *Resolver {
	return &Resolver{Redis: redisClient}
}

// ResolveDomain resolves a domain to a RouteEntry using Redis,
// DNS TXT records and finally the core API.
func (r *Resolver) ResolveDomain(ctx context.Context, domain string) (RouteEntry, error) {
	domain = strings.ToLower(domain)
	key := "routing:" + domain

	// Redis
	val, err := r.Redis.Get(ctx, key).Result()
	if err == nil {
		var route RouteEntry
		if err := json.Unmarshal([]byte(val), &route); err == nil {
			return route, nil
		}
	}

	// TXT lookup
	txts, err := net.LookupTXT(domain)
	if err == nil {
		for _, txt := range txts {
			if strings.HasPrefix(txt, "username_website_") {
				parts := strings.SplitN(txt, "_", 3)
				if len(parts) == 3 {
					username := parts[2]

					route := RouteEntry{
						Target:   "http://127.0.0.1:9999", // can be replaced if username is needed
						Username: username,
					}

					data, _ := json.Marshal(route)
					r.Redis.Set(ctx, key, data, time.Hour)

					return route, nil
				}
			}
		}
	}

	// Request to core API (if there is no TXT and Redis)
	route, err := fetchFromCore(domain)
	if err != nil {
		return RouteEntry{}, err
	}

	// Let's cache
	data, _ := json.Marshal(route)
	r.Redis.Set(ctx, key, data, time.Hour)

	return RouteEntry{}, err
}

// fetchFromCore queries the fallback core API for domain
// resolution when other methods fail.
func fetchFromCore(domain string) (RouteEntry, error) {
	url := fmt.Sprintf("http://127.0.0.1:8080/core/domains/resolve?domain=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return RouteEntry{}, fmt.Errorf("core unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return RouteEntry{}, fmt.Errorf("core API returned status %d", resp.StatusCode)
	}

	var result struct {
		Target   string `json:"target"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RouteEntry{}, fmt.Errorf("invalid response from core: %w", err)
	}

	return RouteEntry{
		Target:   result.Target,
		Username: result.Username,
	}, nil
}
