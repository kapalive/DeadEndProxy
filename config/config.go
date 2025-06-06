package config

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ===== Конфиг для одной локации =====
type LocationConfig struct {
	Path          string `yaml:"path"`
	ProxyPass     string `yaml:"proxy_pass,omitempty"`
	StaticRoot    string `yaml:"static_root,omitempty"`
	IsWebSocket   bool   `yaml:"is_websocket,omitempty"`
	RequireBearer bool   `yaml:"require_bearer,omitempty"`
	Cors          bool   `yaml:"cors,omitempty"`
	Domain        string `yaml:"-"`
}

// ===== Конфиг для сервера =====
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

// ===== Конфиг для YAML =====
type yamlDomain struct {
	Domain          string `yaml:"domain"`
	RedirectToHTTPS bool   `yaml:"redirect_to_https"`
	RedirectTo      string `yaml:"redirect_to,omitempty"`
	SSL             *struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"ssl,omitempty"`
	Routes []LocationConfig `yaml:"routes"`
}

type yamlService struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}

type yamlConfig struct {
	Listen struct {
		HTTP  string `yaml:"http"`
		HTTPS string `yaml:"https"`
	} `yaml:"listen"`
	Domains  []yamlDomain  `yaml:"domains"`
	Services []yamlService `yaml:"services,omitempty"`
}

type Config struct {
	Server ServerConfig
}

// ===== Загрузка конфига =====
func MustLoadConfig(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		panic("config load error: " + err.Error())
	}
	return cfg
}

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
			certMain = main.SSL.Cert
			keyMain = main.SSL.Key
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
			certSecond = second.SSL.Cert
			keySecond = second.SSL.Key
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
			Webroot:       "./webroot", // default, или сделай параметром
			Locations:     allLocations,
		},
	}, nil
}

func parsePort(p string) int {
	return mustAtoi(strings.TrimPrefix(p, ":"))
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
