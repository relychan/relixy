// Package main starts the REST proxy service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/routes/restrouter"
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

	router, shutdown, err := restrouter.SetupRouter(ctx, envVars, ts)
	if err != nil {
		return err
	}

	defer shutdown()

	return gohttps.ListenAndServe(ctx, router, &envVars.Server)
}
