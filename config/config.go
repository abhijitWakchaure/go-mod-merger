package config

import (
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Replace                    map[string]string `json:"replace,omitempty"`
	IgnoreMajorVersionMismatch []string          `json:"ignoreMajorVersionMismatch,omitempty"`
	ForceMaster                []string          `json:"forceMaster"`
}

// Read ...
func Read() *Config {
	return &Config{
		Replace:                    viper.GetStringMapString("replace"),
		IgnoreMajorVersionMismatch: viper.GetStringSlice("ignoreMajorVersionMismatch"),
		ForceMaster:                viper.GetStringSlice("forceMaster"),
	}
}
