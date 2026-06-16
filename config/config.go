// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// EnvNameConfigPath is the constant name of the config path environment variable.
const EnvNameConfigPath = "RELY_CONFIG_PATH"

// RelyDefinitionFileConfig represents the configurations for definition files.
type RelyDefinitionFileConfig struct {
	// List of paths to be included for metadata introspection.
	Include []string `json:"include" yaml:"include"`
	// List of paths to be excluded for metadata introspection.
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

// Validate checks if the definition config is valid.
func (rdc RelyDefinitionFileConfig) Validate() error {
	if len(rdc.Include) == 0 {
		return ErrDefinitionIncludeRequired
	}

	return nil
}

// RelyServerConfig holds information of required configurations to run the Rely API server.
type RelyServerConfig struct {
	// Configurations for the HTTP server.
	Server *gohttps.ServerConfig `json:"server,omitempty" yaml:"server,omitempty"`
	// Configurations for OpenTelemetry exporters.
	Telemetry *gotel.OTLPConfig `json:"telemetry,omitempty" yaml:"telemetry,omitempty"`
	// Configurations for resource definition files.
	Definition RelyDefinitionFileConfig `json:"definition" yaml:"definition"`
}

// LoadServerConfig loads and parses configurations for [RelyServerConfig].
func LoadServerConfig(
	parentContext context.Context,
	configEnvName string,
	serviceName string,
) (*RelyServerConfig, *slog.Logger, error) {
	var result *RelyServerConfig

	var err error

	serverConfigPath := os.Getenv("RELY_CONFIG_PATH")
	if serverConfigPath == "" {
		serverConfigPath = "/etc/rely/config.yaml"
	}

	slog.Info(
		"Loading configurations from file...",
		slog.String("path", serverConfigPath),
	)

	ctx, cancel := context.WithTimeout(parentContext, time.Minute)
	defer cancel()

	result, err = goutils.ReadJSONOrYAMLFile[RelyServerConfig](ctx, serverConfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load RELY_CONFIG_PATH: %w", err)
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
		result.Telemetry.ServiceName = serviceName
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
		pathOrURI, err := goutils.ParseFilePathOrHTTPURL(p)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s: %w", p, err)
		}

		switch {
		case pathOrURI != nil:
			results = append(results, p)
		case filepath.IsAbs(p):
			_, err := filepath.Rel(basePath, p)
			if err != nil {
				return nil, err
			}
		default:
			results = append(results, filepath.Join(basePath, p))
		}
	}

	return results, nil
}
