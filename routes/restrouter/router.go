package restrouter

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/hasura/gotel"
	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/proxyc"
)

func SetupRouter(
	ctx context.Context,
	conf *config.RelixyServerConfig,
	ts *gotel.OTelExporters,
) (*chi.Mux, func(), error) {
	metadata, err := config.LoadMetadata(ctx, conf.Definition)
	if err != nil {
		return nil, nil, err
	}

	httpMetrics, err := gohttpc.NewHTTPClientMetrics(ts.Meter, false)
	if err != nil {
		return nil, nil, err
	}

	gohttpc.SetHTTPClientMetrics(httpMetrics)

	proxyClientOptions := gohttpc.NewClientOptions(
		gohttpc.WithLogger(ts.Logger),
		gohttpc.WithTracer(ts.Tracer),
	)

	middlewares, authManager, err := config.SetupMiddlewares(ctx, &conf.Server, metadata, ts)
	if err != nil {
		return nil, nil, err
	}

	oasResources := metadata.GetOpenAPIResources()

	shutdownFuncs := []func() error{
		authManager.Close,
	}

	shutdown := func() {
		for _, fn := range shutdownFuncs {
			goutils.CatchWarnErrorFunc(fn)
		}
	}

	router := gohttps.NewRouter(&conf.Server, ts.Logger)

	for _, resource := range oasResources {
		proxyClient, err := proxyc.NewProxyClient(ctx, resource, proxyClientOptions)
		if err != nil {
			shutdown()

			return nil, nil, fmt.Errorf("failed to create proxy client %s: %w", resource.Metadata.Name, err)
		}

		var basePath string

		if resource.Definition.Settings != nil && resource.Definition.Settings.BasePath != "" {
			basePath = resource.Definition.Settings.BasePath
		}

		if basePath == "" || basePath[len(basePath)-1] != '/' {
			basePath += "/*"
		} else {
			basePath += "*"
		}

		state := &State{
			ProxyClient: proxyClient,
		}

		shutdownFuncs = append(shutdownFuncs, state.Close)
		router.Handle(
			basePath,
			middlewares.Handler(NewRESTHandler(state)),
		)
	}

	return router, shutdown, nil
}
