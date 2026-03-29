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

package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hasura/gotel"
	"github.com/relychan/relixy/authn"
	"github.com/relychan/rely-auth/auth"
	"github.com/relychan/rely-auth/auth/authmetrics"
	"github.com/relychan/rely-auth/auth/authmode"
)

// SetupMiddlewares sets up default middlewares and the shutdown function for the handler.
func SetupMiddlewares(
	ctx context.Context,
	metadata *RelixyMetadata,
	ts *gotel.OTelExporters,
) (chi.Middlewares, *auth.RelyAuthManager, error) {
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

	var authManager *auth.RelyAuthManager

	if metadata.authResource != nil && len(metadata.authResource.Definition.Modes) > 0 {
		// setup global metrics
		authMetrics, err := authmetrics.NewRelyAuthMetrics(ts.Meter)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to setup auth metrics: %w", err)
		}

		authmetrics.SetRelyAuthMetrics(authMetrics)

		authManager, err = auth.NewRelyAuthManager(
			ctx,
			&auth.RelyAuthConfig{
				Version:    metadata.authResource.Version,
				Definition: metadata.authResource.Definition,
			},
			authmode.WithLogger(ts.Logger),
		)
		if err != nil {
			return nil, nil, err
		}

		middlewares = append(middlewares, authn.AuthMiddleware[map[string]any](authManager))
	}

	return middlewares, authManager, nil
}
