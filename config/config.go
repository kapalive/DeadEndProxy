// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package config provides YAML configuration loading and hot
// reloading logic for DeadEndProxy.
package config

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// ===== Config for one location =====
type LocationConfig struct {
	Path          string `yaml:"path"`
	ProxyPass     string `yaml:"proxy_pass,omitempty"`
	StaticRoot    string `yaml:"static_root,omitempty"`
	IsWebSocket   bool   `yaml:"is_websocket,omitempty"`
	RequireBearer bool   `yaml:"require_bearer,omitempty"`
	Cors          bool   `yaml:"cors,omitempty"`
	Domain        string `yaml:"-"`
}

// ===== Server config =====
type ServerConfig struct {
	HTTPPort      int
	HTTPSPort     int
	DomainMain    string
	DomainSecond  string
	SSLCertMain   string
	SSLKeyMain    string
	SSLCertSecond string
	SSLKeySecond  string
	Webroot       string
	Locations     []LocationConfig
}

// ===== YAML config =====
type yamlDomain struct {
	Domain          string `yaml:"domain"`
	RedirectToHTTPS bool   `yaml:"redirect_to_https"`
	RedirectTo      string `yaml:"redirect_to,omitempty"`
	SSL             *struct {
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"ssl,omitempty"`
	Routes []LocationConfig `yaml:"routes"`
}

type yamlConfig struct {
	Listen struct {
		HTTP  string `yaml:"http"`
		HTTPS string `yaml:"https"`
	} `yaml:"listen"`
	Domains []yamlDomain `yaml:"domains"`
}

type Config struct {
	Server ServerConfig
}

// ===== Loading config =====
// MustLoadConfig loads the YAML configuration and panics
// if an error occurs.
func MustLoadConfig(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		panic("config load error: " + err.Error())
	}
	return cfg
}


// LoadConfig reads the YAML config from disk and converts
// it into the internal Config structure.
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ycfg yamlConfig
	if err := yaml.Unmarshal(raw, &ycfg); err != nil {
		return nil, err
	}

	httpPort := parsePort(ycfg.Listen.HTTP)
	httpsPort := parsePort(ycfg.Listen.HTTPS)

	var allLocations []LocationConfig
	var domainMain, domainSecond, certMain, keyMain, certSecond, keySecond string

	if len(ycfg.Domains) > 0 {
		main := ycfg.Domains[0]
		domainMain = main.Domain
		if main.SSL != nil {
			certMain = main.SSL.CertFile
			keyMain = main.SSL.KeyFile
		}
		for _, route := range main.Routes {
			route.Domain = main.Domain
			allLocations = append(allLocations, route)
		}
	}

	if len(ycfg.Domains) > 1 {
		second := ycfg.Domains[1]
		domainSecond = second.Domain
		if second.SSL != nil {
			certSecond = second.SSL.CertFile
			keySecond = second.SSL.KeyFile
		}
		for _, route := range second.Routes {
			route.Domain = second.Domain
			allLocations = append(allLocations, route)
		}
	}

	return &Config{
		Server: ServerConfig{
			HTTPPort:      httpPort,
			HTTPSPort:     httpsPort,
			DomainMain:    domainMain,
			DomainSecond:  domainSecond,
			SSLCertMain:   certMain,
			SSLKeyMain:    keyMain,
			SSLCertSecond: certSecond,
			SSLKeySecond:  keySecond,
			Webroot:       "./webroot", // default
			Locations:     allLocations,
		},
	}, nil
}

// parsePort trims a leading ':' and converts the string
// representation of a port to an int.
func parsePort(p string) int {
	return mustAtoi(strings.TrimPrefix(p, ":"))
}

// mustAtoi converts string to int and ignores the error.
func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

var (
	currentConfig *Config
	configMu      sync.RWMutex
)

// GetConfig returns the currently loaded configuration.
func GetConfig() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return currentConfig
}

// MustLoadInitial loads the initial config at startup and
// stores it in the global variable.
func MustLoadInitial(path string) {
	cfg, err := LoadConfig(path)
	if err != nil {
		panic("config load error: " + err.Error())
	}
	configMu.Lock()
	currentConfig = cfg
	configMu.Unlock()
}

// WatchAndReload watches the config file and reloads it
// automatically when it changes.
func WatchAndReload(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("fsnotify error:", err)
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("[config] Detected change in config.yaml. Reloading...")
					cfg, err := LoadConfig(path)
					if err != nil {
						log.Println("[config] Reload error:", err)
						continue
					}
					configMu.Lock()
					currentConfig = cfg
					configMu.Unlock()
					log.Println("[config] Reloaded successfully.")
				}
			case err := <-watcher.Errors:
				log.Println("[config] Watcher error:", err)
			}
		}
	}()
	if err := watcher.Add(path); err != nil {
		log.Fatal("watcher.Add:", err)
	}
}
