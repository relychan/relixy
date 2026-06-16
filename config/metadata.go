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
	"errors"
	"path/filepath"
	"strings"

	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema"
	"github.com/relychan/relixy/schema/baseschema"
)

var (
	ErrDefinitionIncludeRequired    = errors.New("require at least 1 included definition")
	ErrAllowOnlyOneRelyAuthResource = errors.New("allow only one RelyAuth resource")
)

// RelyMetadata represents the evaluated rely metadata object.
type RelyMetadata struct {
	openapiResources []*schema.RelyOpenAPIResource
	authResource     *baseschema.RelyAuthResource
}

// LoadMetadata loads metadata resources from definition configurations.
func LoadMetadata(
	ctx context.Context,
	definition RelyDefinitionFileConfig,
) (*RelyMetadata, error) {
	result := &RelyMetadata{}
	includes := []string{}

	for _, include := range definition.Include {
		matches, err := filepath.Glob(include)
		if err != nil {
			return nil, err
		}

	L:
		for _, match := range matches {
			if !isExtensionSupported(match) {
				continue
			}

			for _, exclude := range definition.Exclude {
				shouldExclude, err := filepath.Match(exclude, match)
				if err != nil {
					return nil, err
				}

				if shouldExclude {
					continue L
				}
			}

			includes = append(includes, match)
		}
	}

	if len(includes) == 0 {
		return nil, ErrDefinitionIncludeRequired
	}

	for _, include := range includes {
		resources, err := goutils.ReadMultiFromJSONOrYAMLFile[schema.RelyResource](
			ctx,
			include,
		)
		if err != nil {
			return nil, err
		}

		for _, resource := range resources {
			switch rs := resource.RelyResource.(type) {
			case *baseschema.RelyAuthResource:
				if result.authResource != nil {
					return nil, ErrAllowOnlyOneRelyAuthResource
				}

				result.authResource = rs
			case *schema.RelyOpenAPIResource:
				result.openapiResources = append(result.openapiResources, rs)
			default:
			}
		}
	}

	return result, nil
}

// GetOpenAPIResources returns OpenAPI resources.
func (rm *RelyMetadata) GetOpenAPIResources() []*schema.RelyOpenAPIResource {
	return rm.openapiResources
}

// GetAuthResource returns the RelyAuth resource.
func (rm *RelyMetadata) GetAuthResource() *baseschema.RelyAuthResource {
	return rm.authResource
}

func isExtensionSupported(name string) bool {
	return strings.HasSuffix(name, ".json") ||
		strings.HasSuffix(name, ".yaml") ||
		strings.HasSuffix(name, ".yml")
}
