// © 2023 Devinsidercode CORP. Licensed under the MIT License.
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
	Target   string `json:"target"`   // Прокси-цель
	Username string `json:"username"` // Для аудита
}

type Resolver struct {
	Redis *redis.Client
}

func NewResolver(redisClient *redis.Client) *Resolver {
	return &Resolver{Redis: redisClient}
}

func (r *Resolver) ResolveDomain(ctx context.Context, domain string) (*RouteEntry, error) {
	domain = strings.ToLower(domain)
	key := "routing:" + domain

	// 1️⃣ Redis
	val, err := r.Redis.Get(ctx, key).Result()
	if err == nil {
		var route RouteEntry
		if err := json.Unmarshal([]byte(val), &route); err == nil {
			return &route, nil
		}
	}

	// 2️⃣ TXT lookup
	txts, err := net.LookupTXT(domain)
	if err == nil {
		for _, txt := range txts {
			if strings.HasPrefix(txt, "username_website_") {
				parts := strings.SplitN(txt, "_", 3)
				if len(parts) == 3 {
					username := parts[2]

					route := &RouteEntry{
						Target:   "http://127.0.0.1:9999", // 💥 можно заменить, если username нужен
						Username: username,
					}

					data, _ := json.Marshal(route)
					r.Redis.Set(ctx, key, data, time.Hour)

					return route, nil
				}
			}
		}
	}

	// 3️⃣ Запрос к core API (если нет TXT и Redis)
	route, err := fetchFromCore(domain)
	if err != nil {
		return nil, err
	}

	// 💾 Кешируем
	data, _ := json.Marshal(route)
	r.Redis.Set(ctx, key, data, time.Hour)

	return route, nil
}

func fetchFromCore(domain string) (*RouteEntry, error) {
	url := fmt.Sprintf("http://127.0.0.1:8080/core/domains/resolve?domain=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("core unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("core API returned status %d", resp.StatusCode)
	}

	var result struct {
		Target   string `json:"target"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid response from core: %w", err)
	}

	return &RouteEntry{
		Target:   result.Target,
		Username: result.Username,
	}, nil
}
