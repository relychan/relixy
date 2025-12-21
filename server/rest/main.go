// Package main starts the REST proxy service.
package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relyx/config"
	"github.com/relychan/relyx/routes/rest"
	"github.com/relychan/relyx/types"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	envVars, err := config.LoadServerConfig()
	if err != nil {
		return err
	}

	logger, _, err := otelutils.NewJSONLogger(envVars.Server.LogLevel)
	if err != nil {
		return err
	}

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer stop()

	ts, err := gotel.SetupOTelExporters(ctx, &envVars.Telemetry, types.BuildVersion, logger)
	if err != nil {
		return err
	}

	defer goutils.CatchWarnContextErrorFunc(ts.Shutdown)

	router, shutdown, err := setupRouter(ctx, envVars, ts)
	if err != nil {
		return err
	}

	defer shutdown()

	return gohttps.ListenAndServe(ctx, router, &envVars.Server)
}

func setupRouter(
	ctx context.Context,
	conf *config.RelyXServerConfig,
	ts *gotel.OTelExporters,
) (*chi.Mux, func(), error) {
	state, err := config.NewState(ctx, conf, ts)
	if err != nil {
		return nil, nil, err
	}

	middlewares, shutdown, err := config.SetupMiddlewares(ctx, conf, state, ts)
	if err != nil {
		return nil, nil, err
	}

	basePath := "/*"

	if conf.Router != nil && conf.Router.BasePath != "" {
		basePath = (&url.URL{}).JoinPath(conf.Router.BasePath, "*").RawPath
	}

	router := gohttps.NewRouter(&conf.Server, ts.Logger)
	router.Handle(
		basePath,
		middlewares.Handler(rest.NewRESTHandler(state)),
	)

	return router, shutdown, nil
}
