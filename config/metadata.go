package config

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
)

var (
	ErrDefinitionIncludeRequired    = errors.New("require at least 1 included definition")
	ErrAllowOnlyOneRelyAuthResource = errors.New("allow only one RelyAuth resource")
)

// RelixyMetadata represents the evaluated relixy metadata object.
type RelixyMetadata struct {
	openapiResources []*openapi.RelixyOpenAPIResource
	authResource     *base_schema.RelyAuthResource
}

// LoadMetadata loads metadata resources from definition configurations.
func LoadMetadata(ctx context.Context, definition RelixyDefinitionConfig) (*RelixyMetadata, error) {
	result := &RelixyMetadata{}
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
		resources, err := goutils.ReadMultiFromJSONOrYAMLFile[schema.RelixyResource](
			ctx,
			include,
		)
		if err != nil {
			return nil, err
		}

		for _, resource := range resources {
			switch rs := resource.RelixyResource.(type) {
			case *base_schema.RelyAuthResource:
				if result.authResource != nil {
					return nil, ErrAllowOnlyOneRelyAuthResource
				}

				result.authResource = rs
			case *openapi.RelixyOpenAPIResource:
				result.openapiResources = append(result.openapiResources, rs)
			default:
			}
		}
	}

	return result, nil
}

// GetOpenAPIResources returns OpenAPI resources.
func (rm *RelixyMetadata) GetOpenAPIResources() []*openapi.RelixyOpenAPIResource {
	return rm.openapiResources
}

// GetAuthResource returns the RelyAuth resource.
func (rm *RelixyMetadata) GetAuthResource() *base_schema.RelyAuthResource {
	return rm.authResource
}

func isExtensionSupported(name string) bool {
	return strings.HasSuffix(name, ".json") ||
		strings.HasSuffix(name, ".yaml") ||
		strings.HasSuffix(name, ".yml")
}
