package config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hasura/gotel"
	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/authn"
	"github.com/relychan/relixy/proxyc"
	"github.com/relychan/relixy/schema/openapi"
	"github.com/relychan/relixy/types"
	"github.com/relychan/rely-auth/auth"
	"github.com/relychan/rely-auth/auth/authmetrics"
	"github.com/relychan/rely-auth/auth/authmode"
)

// NewState creates the handler state from config.
func NewState(
	parentContext context.Context,
	conf *RelixyServerConfig,
	ts *gotel.OTelExporters,
) (*types.State, error) {
	ctx, cancel := context.WithTimeout(parentContext, time.Minute)
	defer cancel()

	result, err := goutils.ReadJSONOrYAMLFile[openapi.RelixyOpenAPIResource](ctx, conf.GetConfigPath())
	if err != nil {
		return nil, err
	}

	httpMetrics, err := gohttpc.NewHTTPClientMetrics(ts.Meter, false)
	if err != nil {
		return nil, err
	}

	gohttpc.SetHTTPClientMetrics(httpMetrics)

	proxyClientOptions := gohttpc.NewClientOptions(
		gohttpc.WithLogger(ts.Logger),
		gohttpc.WithTracer(ts.Tracer),
	)

	httpConfig := result.Definition.Settings.HTTP
	if httpConfig == nil {
		httpConfig = new(httpconfig.HTTPClientConfig)
	}

	httpClient, err := httpconfig.NewHTTPClientFromConfig(
		httpConfig,
		proxyClientOptions,
	)
	if err != nil {
		return nil, err
	}

	proxyClientOptions.HTTPClient = httpClient

	proxyClient, err := proxyc.NewProxyClient(&result.Definition, proxyClientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy client: %w", err)
	}

	return &types.State{
		ProxyClient: proxyClient,
	}, nil
}

// SetupMiddlewares sets up default middlewares and the shutdown function for the handler.
func SetupMiddlewares(
	ctx context.Context,
	conf *RelixyServerConfig,
	state *types.State,
	ts *gotel.OTelExporters,
) (chi.Middlewares, func(), error) {
	shutdownFuncs := []func() error{
		state.Close,
	}

	middlewares := chi.Middlewares{
		gotel.NewTracingMiddleware(
			ts,
			gotel.ResponseWriterWrapperFunc(
				func(w http.ResponseWriter, protoMajor int) gotel.WrapResponseWriter {
					return middleware.NewWrapResponseWriter(w, protoMajor)
				},
			),
		),
	}

	if len(conf.Auth.Definition.Modes) > 0 {
		// setup global metrics
		authMetrics, err := authmetrics.NewRelyAuthMetrics(ts.Meter)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to setup auth metrics: %w", err)
		}

		authmetrics.SetRelyAuthMetrics(authMetrics)

		authManager, err := auth.NewRelyAuthManager(
			ctx,
			&conf.Auth,
			authmode.WithLogger(ts.Logger),
		)
		if err != nil {
			goutils.CatchWarnErrorFunc(state.Close)

			return nil, nil, err
		}

		shutdownFuncs = append(shutdownFuncs, authManager.Close)

		middlewares = append(middlewares, authn.AuthMiddleware[map[string]any](authManager))
	}

	shutdown := func() {
		for _, fn := range shutdownFuncs {
			goutils.CatchWarnErrorFunc(fn)
		}
	}

	return middlewares, shutdown, nil
}
