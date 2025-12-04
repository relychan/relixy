package graphqlhandler

import "net/url"

// GraphQLRequestBody represents a request body to a GraphQL server.
type GraphQLRequestBody struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
	Extensions    map[string]any `json:"extensions,omitempty"`
}

type graphqlVariable struct {
	Expression string
	Default    any
}

type requestTemplateData struct {
	Params      map[string]string
	QueryParams url.Values
	Headers     map[string]string
	Body        any
}

// ToMap converts the struct to map.
func (rtd requestTemplateData) ToMap() map[string]any {
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
