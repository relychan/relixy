package ddnrouter

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hasura/gotel"
	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/openapitools/openapiclient"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/schema"
)

// SetupRouter set up the router and states for ddn handlers.
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

	middlewares, authManager, err := config.SetupMiddlewares(ctx, metadata, ts)
	if err != nil {
		return nil, nil, err
	}

	oasResources := metadata.GetOpenAPIResources()

	state := &State{
		ProxyClients: make([]*schema.OpenAPIClient, len(oasResources)),
	}

	shutdown := func() {
		if authManager != nil {
			goutils.CatchWarnErrorFunc(authManager.Close)
		}

		goutils.CatchWarnErrorFunc(state.Close)
	}

	for i, resource := range oasResources {
		proxyClient, err := openapiclient.NewProxyClient(ctx, &resource.Definition)
		if err != nil {
			shutdown()

			return nil, nil, fmt.Errorf(
				"failed to create proxy client %s: %w",
				resource.Metadata.Name,
				err,
			)
		}

		state.ProxyClients[i] = &schema.OpenAPIClient{
			ProxyClient:      proxyClient,
			ResourceMetadata: resource.Metadata,
		}
	}

	router := gohttps.NewRouter(conf.Server, ts.Logger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Handle(
		"/ddn/pre-route",
		middlewares.Handler(NewPreRoutePluginHandler(state)),
	)

	return router, shutdown, nil
}
