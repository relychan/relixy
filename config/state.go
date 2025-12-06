package config

import (
	"fmt"

	"github.com/hasura/gotel"
	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/goutils"
	"github.com/relychan/relyx/proxyc"
	"github.com/relychan/relyx/schema"
	"github.com/relychan/relyx/types"
)

// NewState creates the handler state from config.
func NewState(
	conf *RelyXServerConfig,
	ts *gotel.OTelExporters,
) (*types.State, error) {
	result, err := goutils.ReadJSONOrYAMLFile[schema.RelyProxyAPIDocument](conf.GetConfigPath())
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
