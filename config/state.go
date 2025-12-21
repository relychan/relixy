package config

import (
	"context"
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
	ctx context.Context,
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

	proxyClientOptions := &proxyc.ProxyClientOptions{
		ClientOptions: gohttpc.NewClientOptions(
			gohttpc.WithLogger(ts.Logger),
			gohttpc.WithTracer(ts.Tracer),
		),
	}

	httpConfig := result.Settings.HTTP
	if httpConfig == nil {
		httpConfig = new(httpconfig.HTTPClientConfig)
	}

	httpClient, err := httpconfig.NewHTTPClientFromConfig(
		httpConfig,
		proxyClientOptions.ClientOptions,
	)
	if err != nil {
		return nil, err
	}

	proxyClientOptions.HTTPClient = httpClient

	proxyClient, err := proxyc.NewProxyClient(ctx, result, proxyClientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy client: %w", err)
	}

	return &types.State{
		ProxyClient: proxyClient,
	}, nil
}
