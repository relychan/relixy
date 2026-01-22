// Package main starts the DDN restified endpoint plugin service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/routes/ddn"
	"github.com/relychan/relixy/types"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer stop()

	envVars, err := config.LoadServerConfig(ctx)
	if err != nil {
		return err
	}

	logger, _, err := otelutils.NewJSONLogger(envVars.Server.LogLevel)
	if err != nil {
		return err
	}

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
	conf *config.RelixyServerConfig,
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

	router := gohttps.NewRouter(&conf.Server, ts.Logger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Handle(
		"/ddn/pre-route",
		middlewares.Handler(ddn.NewPreRoutePluginHandler(state)),
	)

	return router, shutdown, nil
}
