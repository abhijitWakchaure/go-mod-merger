package config

import (
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	ForceVersion               map[string]string `json:"forceVersion,omitempty"`
	IgnoreMajorVersionMismatch []string          `json:"ignoreMajorVersionMismatch,omitempty"`
	Remove                     []string          `json:"remove,omitempty"`
	Replace                    map[string]string `json:"replace,omitempty"`
}

// Read ...
func Read() *Config {
	return &Config{
		ForceVersion:               viper.GetStringMapString("forceVersion"),
		IgnoreMajorVersionMismatch: viper.GetStringSlice("ignoreMajorVersionMismatch"),
		Remove:                     viper.GetStringSlice("remove"),
		Replace:                    viper.GetStringMapString("replace"),
	}
}
