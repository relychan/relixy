// Package config defines configurations to start the server.
package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hasura/gotel"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
)

// RelixyDefinitionConfig represents the configurations for relixy definitions.
type RelixyDefinitionConfig struct {
	// List of paths to be included for metadata introspection.
	Include []string `json:"include" yaml:"include"`
	// List of paths to be excluded for metadata introspection.
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

// Validate checks if the definition config is valid.
func (rdc RelixyDefinitionConfig) Validate() error {
	if len(rdc.Include) == 0 {
		return ErrDefinitionIncludeRequired
	}

	return nil
}

// RelixyServerConfig holds information of required configurations to run the relixy server.
type RelixyServerConfig struct {
	Server     gohttps.ServerConfig   `json:"server" yaml:"server"`
	Telemetry  gotel.OTLPConfig       `json:"telemetry" yaml:"telemetry"`
	Definition RelixyDefinitionConfig `json:"definition" yaml:"definition"`
}

// LoadServerConfig loads and parses configurations for [RelyAuthServerConfig].
func LoadServerConfig(parentContext context.Context) (*RelixyServerConfig, error) {
	var result *RelixyServerConfig

	var err error

	serverConfigPath := os.Getenv("RELIXY_CONFIG_PATH")
	if serverConfigPath != "" {
		ctx, cancel := context.WithTimeout(parentContext, time.Minute)
		defer cancel()

		result, err = goutils.ReadJSONOrYAMLFile[RelixyServerConfig](ctx, serverConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load RELIXY_CONFIG_PATH: %w", err)
		}
	} else {
		result = &RelixyServerConfig{}
	}

	err = env.Parse(result)
	if err != nil {
		return result, fmt.Errorf("failed to load environment variables for server config: %w", err)
	}

	err = result.Definition.Validate()
	if err != nil {
		return nil, err
	}

	if result.Telemetry.ServiceName == "" {
		result.Telemetry.ServiceName = "relixy"
	}

	return result, nil
}
