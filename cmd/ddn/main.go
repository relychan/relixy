// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main starts the DDN restified endpoint plugin service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/hasura/gotel"
	"github.com/relychan/gohttps"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/routes/ddnrouter"
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

	envVars, logger, err := config.LoadServerConfig(ctx)
	if err != nil {
		return err
	}

	ts, err := gotel.SetupOTelExporters(ctx, envVars.Telemetry, config.BuildVersion, logger)
	if err != nil {
		return err
	}

	defer goutils.CatchWarnContextErrorFunc(ts.Shutdown)

	router, shutdown, err := ddnrouter.SetupRouter(ctx, envVars, ts)
	if err != nil {
		return err
	}

	defer shutdown()

	return gohttps.ListenAndServe(ctx, router, envVars.Server)
}
