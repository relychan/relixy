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

// RelyXRouterConfig holds configurations of the rest handler.
type RelyXRouterConfig struct {
	// Set the base path for all API handlers.
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty" env:"RELYX_ROUTE_BASE_PATH"`
}

// RelyXServerConfig holds information of required configurations to run the relyx server.
type RelyXServerConfig struct {
	Server     gohttps.ServerConfig `json:"server" yaml:"server"`
	Router     RelyXRouterConfig    `json:"router,omitempty" yaml:"router"`
	Telemetry  gotel.OTLPConfig     `json:"telemetry" yaml:"telemetry"`
	Auth       auth.RelyAuthConfig  `json:"auth" yaml:"auth"`
	ConfigPath string               `json:"configPath" yaml:"configPath" env:"RELYX_CONFIG_PATH"`
}

// GetConfigPath returns the auth config path.
func (rlsc RelyXServerConfig) GetConfigPath() string {
	if rlsc.ConfigPath != "" {
		return rlsc.ConfigPath
	}

	return "/etc/relyx/config.yaml"
}

// LoadServerConfig loads and parses configurations for [RelyAuthServerConfig].
func LoadServerConfig() (*RelyXServerConfig, error) {
	var result *RelyXServerConfig

	var err error

	serverConfigPath := os.Getenv("RELYX_SERVER_CONFIG_PATH")
	if serverConfigPath != "" {
		result, err = goutils.ReadJSONOrYAMLFile[RelyXServerConfig](serverConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load RELY_AUTH_SERVER_CONFIG_PATH: %w", err)
		}
	} else {
		result = &RelyXServerConfig{}
	}

	err = env.Parse(result)
	if err != nil {
		return result, fmt.Errorf("failed to load environment variables for server config: %w", err)
	}

	if result.Telemetry.ServiceName == "" {
		result.Telemetry.ServiceName = "relyx"
	}

	return result, nil
}
