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

	state, err := config.NewState(envVars, ts)
	if err != nil {
		return err
	}

	defer goutils.CatchWarnErrorFunc(state.Close)

	router := setupRouter(state, envVars, ts)

	return gohttps.ListenAndServe(ctx, router, &envVars.Server)
}

func setupRouter(
	state *types.State,
	envVars *config.RelyXServerConfig,
	ts *gotel.OTelExporters,
) *chi.Mux {
	router := gohttps.NewRouter(&envVars.Server, ts.Logger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Handle(
		"/ddn/pre-route",
		chi.Chain(
			gotel.NewTracingMiddleware(
				ts,
				gotel.ResponseWriterWrapperFunc(
					func(w http.ResponseWriter, protoMajor int) gotel.WrapResponseWriter {
						return middleware.NewWrapResponseWriter(w, protoMajor)
					},
				),
			), authn.AuthMiddleware[map[string]any](nil),
		).Handler(ddn.NewPreRoutePluginHandler(state)),
	)

	return router
}
