package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/hasura/gotel"
	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relyx/proxyc"
	"github.com/relychan/relyx/schema"
	"github.com/relychan/relyx/types"
)

// Environment holds information of required environment variables.
type Environment struct {
	Server    gohttps.ServerConfig
	Telemetry gotel.OTLPConfig

	ConfigPath string `env:"CONFIG_PATH" envDefault:"/etc/relyx/config.yaml"`
}

// GetEnvironment loads and parses environment variables.
func GetEnvironment() Environment {
	result, err := env.ParseAs[Environment]()
	if err != nil {
		log.Fatalf("failed to parse environment variables: %s", err) //nolint:revive
	}

	if result.ConfigPath == "" {
		log.Fatal("config path is required") //nolint:revive
	}

	if result.Telemetry.ServiceName == "" {
		result.Telemetry.ServiceName = "relyx-ddn"
	}

	return result
}

// NewState creates the handler state from config.
func NewState(
	environment *Environment,
	ts *gotel.OTelExporters,
) (*types.State, error) {
	result, err := goutils.ReadJSONOrYAMLFile[schema.RelyProxyAPIDocument](environment.ConfigPath)
	if err != nil {
		return nil, err
	}

	httpMetrics, err := gohttpc.NewHTTPClientMetrics(ts.Meter, false)
	if err != nil {
		return nil, err
	}

	gohttpc.SetHTTPClientMetrics(httpMetrics)

	clientOptions := gohttpc.NewClientOptions(
		gohttpc.WithLogger(ts.Logger),
		gohttpc.WithTracer(ts.Tracer),
	)

	httpConfig := result.Settings.HTTP
	if httpConfig == nil {
		httpConfig = new(httpconfig.HTTPClientConfig)
	}

	httpClient, err := httpconfig.NewHTTPClientFromConfig(
		*httpConfig,
		clientOptions,
	)
	if err != nil {
		return nil, err
	}

	clientOptions.HTTPClient = httpClient

	proxyClient, err := proxyc.NewProxyClient(result, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy client: %w", err)
	}

	return &types.State{
		ProxyClient: proxyClient,
	}, nil
}
