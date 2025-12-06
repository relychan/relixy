// Package main starts the DDN restified endpoint plugin service.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/rely-auth/auth"
	"github.com/relychan/relyx/authn"
	"github.com/relychan/relyx/config"
	"github.com/relychan/relyx/routes/ddn"
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

	router, shutdown, err := setupRouter(envVars, ts)
	if err != nil {
		return err
	}

	defer shutdown()

	return gohttps.ListenAndServe(ctx, router, &envVars.Server)
}

func setupRouter(
	conf *config.RelyXServerConfig,
	ts *gotel.OTelExporters,
) (*chi.Mux, func(), error) {
	state, err := config.NewState(conf, ts)
	if err != nil {
		return nil, nil, err
	}

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

	if len(conf.Auth.Definitions) > 0 {
		authManager, err := auth.NewRelyAuthManager(
			&conf.Auth,
			auth.WithLogger(ts.Logger),
			auth.WithMeter(ts.Meter),
		)
		if err != nil {
			goutils.CatchWarnErrorFunc(state.Close)

			return nil, nil, err
		}

		shutdownFuncs = append(shutdownFuncs, authManager.Close)

		middlewares = append(middlewares, authn.AuthMiddleware[map[string]any](authManager))
	}

	router := gohttps.NewRouter(&conf.Server, ts.Logger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Handle(
		"/ddn/pre-route",
		middlewares.Handler(ddn.NewPreRoutePluginHandler(state)),
	)

	shutdown := func() {
		for _, fn := range shutdownFuncs {
			goutils.CatchWarnErrorFunc(fn)
		}
	}

	return router, shutdown, nil
}
