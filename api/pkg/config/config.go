package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type OIDCConfig struct {
	DisplayName  string   `yaml:"display_name"`
	Issuer       string   `yaml:"issuer"`
	ClientID     string   `yaml:"clientId"`
	ClientSecret string   `yaml:"clientSecret"`
	Scopes       []string `yaml:"scopes"`
}

type SecurityConfig struct {
	SigninigKey   string                `yaml:"signingKey"`
	AdminPassword string                `yaml:"adminPassword"`
	OIDC          map[string]OIDCConfig `yaml:"oidc"`
}

type HTTPConfig struct {
	Listen string `yaml:"listen"`
}

type StorageConfig struct {
	Driver string `yaml:"driver"`
	URL    string `yaml:"url"`
}

type Config struct {
	Debug    bool           `yaml:"debug"`
	Storage  StorageConfig  `yaml:"storage"`
	HTTP     HTTPConfig     `yaml:"http"`
	Security SecurityConfig `yaml:"security"`
}

func LoadFromFile(pth string) (*Config, error) {
	b, err := os.ReadFile(pth)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(b, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
