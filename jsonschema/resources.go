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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"

	"github.com/relychan/goutils"
	"github.com/relychan/jsonschema"
	"github.com/relychan/relixy/schema"
)

func genOpenAPIResourceSchema() (*jsonschema.Schema, error) {
	r := new(jsonschema.Reflector)

	err := r.AddGoComments(
		"github.com/relychan/relixy/schema",
		"../schema",
		jsonschema.WithFullComment(),
	)
	if err != nil {
		return nil, err
	}

	reflectSchema := r.Reflect(schema.RelixyOpenAPIResource{})

	remoteSchemas, err := downloadRemoteSchemas()
	if err != nil {
		return nil, err
	}

	for _, rs := range remoteSchemas {
		maps.Copy(reflectSchema.Definitions, rs.Definitions)
	}

	// delete unused definitions
	delete(reflectSchema.Definitions, "Contact")
	delete(reflectSchema.Definitions, "Components")
	delete(reflectSchema.Definitions, "ExternalDoc")
	delete(reflectSchema.Definitions, "Tag")
	delete(reflectSchema.Definitions, "SecurityRequirement")
	delete(reflectSchema.Definitions, "Server")
	delete(reflectSchema.Definitions, "Paths")
	delete(reflectSchema.Definitions, "Info")
	delete(reflectSchema.Definitions, "License")

	for key := range reflectSchema.Definitions {
		if strings.HasPrefix(key, "Map[") {
			delete(reflectSchema.Definitions, key)
		}
	}

	return reflectSchema, nil
}

func downloadRemoteSchemas() ([]*jsonschema.Schema, error) {
	fileURLs := []string{
		"https://raw.githubusercontent.com/relychan/rely-auth/refs/heads/main/jsonschema/auth.schema.json",
		"https://raw.githubusercontent.com/relychan/openapitools/refs/heads/main/jsonschema/openapitools.schema.json",
	}

	results := make([]*jsonschema.Schema, 0, len(fileURLs))

	for _, fileURL := range fileURLs {
		rawResp, err := http.Get(fileURL) //nolint:bodyclose,noctx,gosec
		if err != nil {
			return nil, fmt.Errorf("failed to download file %s: %w", fileURL, err)
		}

		if rawResp != nil && rawResp.Body != nil {
			defer goutils.CatchWarnErrorFunc(rawResp.Body.Close) //nolint:revive
		}

		if rawResp.StatusCode != http.StatusOK {
			rawBody, err := io.ReadAll(rawResp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to download %s schema: %s", fileURL, rawResp.Status) //nolint
			}

			return nil, fmt.Errorf("failed to download %s schema: %s", fileURL, string(rawBody)) //nolint
		}

		jsonSchema := new(jsonschema.Schema)

		err = json.NewDecoder(rawResp.Body).Decode(jsonSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to decode gohttpc schema: %w", err)
		}

		results = append(results, jsonSchema)
	}

	return results, nil
}
