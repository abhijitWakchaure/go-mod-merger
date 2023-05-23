package config

import (
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Replace map[string]string `json:"replace,omitempty"`
}

// Read ...
func Read() *Config {
	c := &Config{}
	m := viper.GetStringMapString("replace")
	c.Replace = m
	return c
}
