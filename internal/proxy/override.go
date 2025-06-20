// © 2023 Devinsidercode CORP. Licensed under the MIT License.
package proxy

import (
	"fmt"
	"strconv"

	"DeadEndProxy/config"
)

type ConfigOverride struct {
	HTTPPort  string
	HTTPSPort string
}

func (o *ConfigOverride) Apply(cfg *config.Config) {
	if o.HTTPPort != "" {
		port := mustInt(o.HTTPPort)
		fmt.Printf("🔧 Overriding HTTP port: %d -> %d\n", cfg.Server.HTTPPort, port)
		cfg.Server.HTTPPort = port
	}

	if o.HTTPSPort != "" {
		port := mustInt(o.HTTPSPort)
		fmt.Printf("🔧 Overriding HTTPS port: %d -> %d\n", cfg.Server.HTTPSPort, port)
		cfg.Server.HTTPSPort = port
	}
}

func mustInt(val string) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		panic("Invalid port override: " + val)
	}
	return n
}
