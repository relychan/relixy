package proxyhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform"
	"github.com/relychan/goutils"
	"github.com/relychan/goutils/httpheader"
)

var (
	errOAuth2ClientCredentialsRequired = errors.New("clientId and clientSecret must not be empty")
	errOAuth2TokenURLRequired          = errors.New(
		"tokenUrl is required in the OAuth2 Client Credentials flow",
	)
)

// ProxyActionType represents enums of proxy types.
type ProxyActionType string

// InsertRouteOptions represents options for inserting routes.
type InsertRouteOptions struct {
	GetEnv goenvconf.GetEnvFunc
}

// APIKeyCredentials holds apiKey credentials of the security scheme.
type APIKeyCredentials struct {
	APIKey *goenvconf.EnvString `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
}

// BasicCredentials holds basic credentials of the security scheme.
type BasicCredentials struct {
	Username *goenvconf.EnvString `json:"username,omitempty" yaml:"username,omitempty"`
	Password *goenvconf.EnvString `json:"password,omitempty" yaml:"password,omitempty"`
}

// OAuth2Credentials holds OAuth2 credentials of the security scheme.
type OAuth2Credentials struct {
	ClientID     *goenvconf.EnvString `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret *goenvconf.EnvString `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	// Optional query parameters for the token and refresh URLs.
	EndpointParams map[string]goenvconf.EnvString `json:"endpointParams,omitempty" yaml:"endpointParams,omitempty"`
}

// RelixyCustomResponseConfig represents configurations for the proxy response.
type RelixyCustomResponseConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int `json:"httpErrorCode,omitempty" yaml:"httpErrorCode,omitempty" jsonschema:"minimum=400,maximum=599,default=400"`
	// Configurations for transforming response data.
	Transform *gotransform.TemplateTransformerConfig `json:"transform,omitempty" yaml:"transform,omitempty"`
}

// IsZero checks if the configuration is empty.
func (conf RelixyCustomResponseConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Transform == nil || conf.Transform.IsZero())
}

// RelixyCustomResponse represents configurations for the proxy response.
type RelixyCustomResponse struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int
	// Configurations for transforming response body data.
	Body gotransform.TemplateTransformer
}

// NewRelixyCustomResponse creates a [RelixyCustomResponse] from raw configurations.
func NewRelixyCustomResponse(
	config *RelixyCustomResponseConfig,
	getEnv goenvconf.GetEnvFunc,
) (*RelixyCustomResponse, error) {
	if config == nil || config.IsZero() {
		return nil, nil
	}

	result := &RelixyCustomResponse{
		HTTPErrorCode: config.HTTPErrorCode,
	}

	if config.Transform != nil {
		transformer, err := gotransform.NewTransformerFromConfig("", *config.Transform, getEnv)
		if err != nil {
			return result, err
		}

		result.Body = transformer
	}

	return result, nil
}

// IsZero checks if the configuration is empty.
func (conf RelixyCustomResponse) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Body == nil || conf.Body.IsZero())
}

// RequestTemplateData represents the request data for template transformation.
type RequestTemplateData struct {
	Params      map[string]string
	QueryParams url.Values
	Headers     map[string]string
	Body        any
}

// NewRequestTemplateData creates a new [RequestTemplateData] from the HTTP request to a map for request transformation.
func NewRequestTemplateData(
	request *http.Request,
	contentType string,
	paramValues map[string]string,
) (*RequestTemplateData, bool, error) {
	requestHeaders := map[string]string{}

	for key, header := range request.Header {
		if len(header) == 0 {
			continue
		}

		requestHeaders[strings.ToLower(key)] = header[0]
	}

	requestData := &RequestTemplateData{
		Params:      paramValues,
		QueryParams: request.URL.Query(),
		Headers:     requestHeaders,
	}

	if request.Body == nil || request.Body == http.NoBody {
		return requestData, true, nil
	}

	switch {
	case strings.HasPrefix(contentType, httpheader.ContentTypeJSON):
		defer goutils.CatchWarnErrorFunc(request.Body.Close)

		var body any

		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			return nil, true, err
		}

		requestData.Body = body
	default:
		// skip other content types.
	}

	return requestData, false, nil
}

// ToMap converts the struct to map.
func (rtd RequestTemplateData) ToMap() map[string]any {
	result := map[string]any{
		"param":   rtd.Params,
		"query":   rtd.QueryParams,
		"headers": rtd.Headers,
	}

	if rtd.Body != nil {
		result["body"] = rtd.Body
	}

	return result
}
