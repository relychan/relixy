// Package config defines configurations to start the server.
package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/hasura/gotel"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/rely-auth/auth"
)

// RelixyServerConfig holds information of required configurations to run the relixy server.
type RelixyServerConfig struct {
	Server     gohttps.ServerConfig `json:"server" yaml:"server"`
	Telemetry  gotel.OTLPConfig     `json:"telemetry" yaml:"telemetry"`
	Auth       auth.RelyAuthConfig  `json:"auth" yaml:"auth"`
	ConfigPath string               `json:"configPath" yaml:"configPath" env:"RELIXY_CONFIG_PATH"`
}

// GetConfigPath returns the auth config path.
func (rlsc RelixyServerConfig) GetConfigPath() string {
	if rlsc.ConfigPath != "" {
		return rlsc.ConfigPath
	}

	return "/etc/relixy/config.yaml"
}

// LoadServerConfig loads and parses configurations for [RelyAuthServerConfig].
func LoadServerConfig() (*RelixyServerConfig, error) {
	var result *RelixyServerConfig

	var err error

	serverConfigPath := os.Getenv("RELIXY_SERVER_CONFIG_PATH")
	if serverConfigPath != "" {
		result, err = goutils.ReadJSONOrYAMLFile[RelixyServerConfig](serverConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load RELIXY_SERVER_CONFIG_PATH: %w", err)
		}
	} else {
		result = &RelixyServerConfig{}
	}

	err = env.Parse(result)
	if err != nil {
		return result, fmt.Errorf("failed to load environment variables for server config: %w", err)
	}

	if result.Telemetry.ServiceName == "" {
		result.Telemetry.ServiceName = "relixy"
	}

	return result, nil
}
