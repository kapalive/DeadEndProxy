package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port     int      `yaml:"port"`
		LogLevel string   `yaml:"log_level"`
		Hosts    []string `yaml:"allowed_hosts"`
	} `yaml:"server"`
	Proxy struct {
		Routes []struct {
			Path    string   `yaml:"path"`
			Backend string   `yaml:"backend"`
			Methods []string `yaml:"methods"`
		} `yaml:"routes"`
	} `yaml:"proxy"`
	Security struct {
		BlockedIPs        []string `yaml:"blocked_ips"`
		BlockedUserAgents []string `yaml:"blocked_user_agents"`
		EnableWAF         bool     `yaml:"enable_waf"`
		RateLimit         int      `yaml:"ratelimit"`
	} `yaml:"security"`
}

// LoadConfig reads config.yaml and parses it into a structure
func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
