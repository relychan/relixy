package resthandler

import (
	"fmt"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
)

// ProxyActionTypeREST represents a constant value for REST proxy action.
const ProxyActionTypeREST proxyhandler.ProxyActionType = "rest"

// RelixyRESTActionConfig represents a proxy action config for REST operation.
type RelixyRESTActionConfig struct {
	// Type of the proxy action which is always rest.
	Type proxyhandler.ProxyActionType `json:"type" yaml:"type" jsonschema:"enum=rest"`
	// Configurations for the REST proxy request.
	Request *RelixyRESTRequestConfig `json:"request,omitempty" yaml:"request,omitempty"`
	// Configurations for evaluating REST responses.
	Response *RelixyCustomRESTResponseConfig `json:"response,omitempty" yaml:"response,omitempty"`
}

// RelixyCustomRESTResponseConfig represents configurations for the proxy response.
type RelixyCustomRESTResponseConfig struct {
	// Configurations for transforming response data.
	Body *gotransform.TemplateTransformerConfig `json:"body,omitempty" yaml:"body,omitempty"`
}

// IsZero checks if the configuration is empty.
func (conf RelixyCustomRESTResponseConfig) IsZero() bool {
	return conf.Body == nil || conf.Body.IsZero()
}

type customRESTResponse struct {
	// Configurations for transforming response body data.
	Body gotransform.TemplateTransformer
}

// newCustomRESTResponse creates a [RelixyCustomResponse] from raw configurations.
func newCustomRESTResponse(
	config *RelixyCustomRESTResponseConfig,
	getEnv goenvconf.GetEnvFunc,
) (*customRESTResponse, error) {
	if config == nil || config.IsZero() {
		return nil, nil
	}

	transformer, err := gotransform.NewTransformerFromConfig("", *config.Body, getEnv)
	if err != nil {
		return nil, err
	}

	result := &customRESTResponse{
		Body: transformer,
	}

	return result, nil
}

// IsZero checks if the configuration is empty.
func (conf customRESTResponse) IsZero() bool {
	return conf.Body == nil || conf.Body.IsZero()
}

// RelixyRESTRequestConfig represents configurations for the proxy request.
type RelixyRESTRequestConfig struct {
	// Overrides the request path. Use the original request path if empty.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	// The configuration to transform request headers.
	Headers map[string]jmes.FieldMappingEntryStringConfig `json:"headers,omitempty" yaml:"headers,omitempty"`
	// The configuration to transform request body.
	Body *gotransform.TemplateTransformerConfig `json:"body,omitempty" yaml:"body,omitempty"`
}

// IsZero checks if the configuration is empty.
func (rr RelixyRESTRequestConfig) IsZero() bool {
	return rr.Path == "" &&
		len(rr.Headers) == 0 &&
		(rr.Body == nil || rr.Body.IsZero())
}

type customRESTRequest struct {
	Path    string
	Headers map[string]jmes.FieldMappingEntryString
	Body    gotransform.TemplateTransformer
}

// IsZero checks if the configuration is empty.
func (rr customRESTRequest) IsZero() bool {
	return rr.Path == "" &&
		len(rr.Headers) == 0 &&
		(rr.Body == nil || rr.Body.IsZero())
}

func newCustomRESTRequestFromConfig(
	conf *RelixyRESTRequestConfig,
	getEnvFunc goenvconf.GetEnvFunc,
) (*customRESTRequest, error) {
	if conf == nil || conf.IsZero() {
		return nil, nil
	}

	result := &customRESTRequest{
		Path: conf.Path,
	}

	if len(conf.Headers) > 0 {
		headers, err := jmes.EvaluateObjectFieldMappingStringEntries(
			conf.Headers, getEnvFunc)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize custom request headers: %w", err)
		}

		result.Headers = headers
	}

	if conf.Body != nil {
		customBody, err := gotransform.NewTransformerFromConfig("", *conf.Body, getEnvFunc)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize custom request body: %w", err)
		}

		result.Body = customBody
	}

	return result, nil
}
