// Package config defines configurations to start the server.
package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
)

// BuildVersion is set when building the binary.
var BuildVersion = "dev"

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
	// Configurations for the HTTP server.
	Server *gohttps.ServerConfig `json:"server,omitempty" yaml:"server,omitempty"`
	// Configurations for OpenTelemetry exporters.
	Telemetry *gotel.OTLPConfig `json:"telemetry,omitempty" yaml:"telemetry,omitempty"`
	// Configurations for resource definition files.
	Definition RelixyDefinitionConfig `json:"definition" yaml:"definition"`
}

// LoadServerConfig loads and parses configurations for [RelixyServerConfig].
func LoadServerConfig(parentContext context.Context) (*RelixyServerConfig, *slog.Logger, error) {
	var result *RelixyServerConfig

	var err error

	serverConfigPath := os.Getenv("RELIXY_CONFIG_PATH")
	if serverConfigPath != "" {
		slog.Debug( //nolint:gosec
			"Loading configurations from file...",
			slog.String("path", serverConfigPath),
		)

		ctx, cancel := context.WithTimeout(parentContext, time.Minute)
		defer cancel()

		result, err = goutils.ReadJSONOrYAMLFile[RelixyServerConfig](ctx, serverConfigPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load RELIXY_CONFIG_PATH: %w", err)
		}
	} else {
		result = &RelixyServerConfig{}
	}

	if result.Server == nil {
		result.Server = new(gohttps.ServerConfig)
	}

	if result.Telemetry == nil {
		result.Telemetry = new(gotel.OTLPConfig)
	}

	err = env.Parse(result)
	if err != nil {
		return result,
			nil,
			fmt.Errorf("failed to load environment variables for server config: %w", err)
	}

	logger, _, err := otelutils.NewJSONLogger(result.Server.GetLogLevel())
	if err != nil {
		return nil, nil, err
	}

	logger.Debug("Loaded configurations", slog.Any("config", result))

	err = result.Definition.Validate()
	if err != nil {
		return nil, nil, err
	}

	if result.Telemetry.ServiceName == "" {
		result.Telemetry.ServiceName = "relixy"
	}

	if serverConfigPath != "" && serverConfigPath != "." {
		basePath := filepath.Dir(serverConfigPath)

		include, err := resolveDefinitionPaths(basePath, result.Definition.Include)
		if err != nil {
			return nil, nil, err
		}

		result.Definition.Include = include

		exclude, err := resolveDefinitionPaths(basePath, result.Definition.Exclude)
		if err != nil {
			return nil, nil, err
		}

		result.Definition.Exclude = exclude
	}

	return result, logger, nil
}

func resolveDefinitionPaths(basePath string, paths []string) ([]string, error) {
	results := make([]string, 0, len(paths))

	for _, p := range paths {
		pathOrURI, err := goutils.ParsePathOrHTTPURL(p)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s: %w", p, err)
		}

		if pathOrURI.Scheme != "" || filepath.IsAbs(p) {
			results = append(results, p)
		} else {
			results = append(results, filepath.Join(basePath, p))
		}
	}

	return results, nil
}
