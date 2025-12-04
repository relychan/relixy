package main

import (
	"github.com/invopop/jsonschema"
	"github.com/relychan/gohttpc/authc/authscheme"
	"github.com/relychan/goutils"
	"github.com/relychan/relyx/schema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func newRelyProxyGraphQLRequestConfigSchema() *jsonschema.Schema {
	graphqlProps := orderedmap.New[string, *jsonschema.Schema]()

	graphqlProps.Set("query", &jsonschema.Schema{
		Description: "GraphQL query string to send",
		Type:        "string",
	})
	graphqlProps.Set("variables", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/GraphQLVariableDefinition",
		},
	})
	graphqlProps.Set("extensions", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/GraphQLVariableDefinition",
		},
	})

	return &jsonschema.Schema{
		Type:       "object",
		Properties: graphqlProps,
		Required:   []string{"query"},
	}
}

func newRelyProxyActionSchema() *jsonschema.Schema {
	restSchema := orderedmap.New[string, *jsonschema.Schema]()
	restSchema.Set("type", &jsonschema.Schema{
		Type:        "string",
		Description: "Type of the proxy action",
		Enum:        []any{schema.ProxyTypeREST},
	})
	restSchema.Set("path", &jsonschema.Schema{
		Description: "Overrides the request path. Use the original request path if empty",
		Type:        "string",
	})

	graphqlSchema := orderedmap.New[string, *jsonschema.Schema]()
	graphqlSchema.Set("type", &jsonschema.Schema{
		Type:        "string",
		Description: "Type of the proxy action",
		Enum:        []any{schema.ProxyTypeGraphQL},
	})
	graphqlSchema.Set("request", &jsonschema.Schema{
		Description: "Configuration for the GraphQL request",
		Ref:         "#/$defs/RelyProxyGraphQLRequestConfig",
	})
	graphqlSchema.Set("response", &jsonschema.Schema{
		Description: "Configuration for the GraphQL response",
		Ref:         "#/$defs/RelyProxyGraphQLResponseConfig",
	})

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "object",
				Description: "Proxy configuration to the remote REST service",
				Required:    []string{"type"},
				Properties:  restSchema,
			},
			{
				Type:        "object",
				Description: "Configurations for proxying request to the remote GraphQL server",
				Properties:  graphqlSchema,
				Required:    []string{"type", "request"},
			},
		},
	}
}

func newRelyProxySecuritySchemeSchema() *jsonschema.Schema {
	envStringRef := &jsonschema.Schema{
		Ref: "#/$defs/EnvString",
	}
	apiKeySchema := orderedmap.New[string, *jsonschema.Schema]()
	apiKeySchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.APIKeyScheme},
	})
	apiKeySchema.Set("value", envStringRef)
	apiKeySchema.Set("in", &jsonschema.Schema{
		Type: "string",
		Enum: goutils.ToAnySlice(authscheme.GetSupportedAuthLocations()),
	})
	apiKeySchema.Set("name", &jsonschema.Schema{
		Type: "string",
	})

	httpAuthSchema := orderedmap.New[string, *jsonschema.Schema]()
	httpAuthSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.HTTPAuthScheme},
	})
	httpAuthSchema.Set("value", envStringRef)
	httpAuthSchema.Set("header", &jsonschema.Schema{
		Type: "string",
	})
	httpAuthSchema.Set("scheme", &jsonschema.Schema{
		Type: "string",
	})

	basicAuthSchema := orderedmap.New[string, *jsonschema.Schema]()
	basicAuthSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.BasicAuthScheme},
	})
	basicAuthSchema.Set("username", envStringRef)
	basicAuthSchema.Set("password", envStringRef)
	httpAuthSchema.Set("header", &jsonschema.Schema{
		Description: "Request contains a header field in the form of Authorization: Basic [credentials]",
		OneOf: []*jsonschema.Schema{
			{Type: "string"},
			{Type: "null"},
		},
	})

	oidcSchema := orderedmap.New[string, *jsonschema.Schema]()
	oidcSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.OpenIDConnectScheme},
	})
	oidcSchema.Set("openIdConnectUrl", &jsonschema.Schema{
		Type: "string",
	})

	cookieSchema := orderedmap.New[string, *jsonschema.Schema]()
	cookieSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.CookieAuthScheme},
	})

	mutualTLSSchema := orderedmap.New[string, *jsonschema.Schema]()
	mutualTLSSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.MutualTLSScheme},
	})

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:       "object",
				Required:   []string{"type", "value", "in", "name"},
				Properties: apiKeySchema,
			},
			{
				Type:       "object",
				Properties: basicAuthSchema,
				Required:   []string{"type", "username", "password"},
			},
			{
				Type:       "object",
				Properties: httpAuthSchema,
				Required:   []string{"type", "value", "header", "scheme"},
			},
			newOAuth2ConfigSchema(),
			{
				Type:       "object",
				Properties: oidcSchema,
				Required:   []string{"type", "openIdConnectUrl"},
			},
			{
				Type:       "object",
				Properties: cookieSchema,
				Required:   []string{"type"},
			},
			{
				Type:       "object",
				Properties: mutualTLSSchema,
				Required:   []string{"type"},
			},
		},
	}
}

func newOAuth2ConfigSchema() *jsonschema.Schema {
	oauth2Schema := orderedmap.New[string, *jsonschema.Schema]()
	oauth2Schema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{schema.OAuth2Scheme},
	})

	oauthFlowRef := &jsonschema.Schema{
		Ref: "#/$defs/OAuthFlow",
	}
	implicitFlow := orderedmap.New[string, *jsonschema.Schema]()
	implicitFlow.Set(string(schema.ImplicitFlow), oauthFlowRef)

	passwordFlow := orderedmap.New[string, *jsonschema.Schema]()
	passwordFlow.Set(string(schema.PasswordFlow), oauthFlowRef)

	ccFlow := orderedmap.New[string, *jsonschema.Schema]()
	ccFlow.Set(string(schema.ClientCredentialsFlow), oauthFlowRef)

	acFlow := orderedmap.New[string, *jsonschema.Schema]()
	acFlow.Set(string(schema.AuthorizationCodeFlow), oauthFlowRef)

	oauth2Schema.Set("flows", &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:       "object",
				Required:   []string{string(schema.PasswordFlow)},
				Properties: passwordFlow,
			},
			{
				Type:       "object",
				Required:   []string{string(schema.ImplicitFlow)},
				Properties: implicitFlow,
			},
			{
				Type:       "object",
				Required:   []string{string(schema.ClientCredentialsFlow)},
				Properties: ccFlow,
			},
			{
				Type:       "object",
				Required:   []string{string(schema.AuthorizationCodeFlow)},
				Properties: acFlow,
			},
		},
	})

	return &jsonschema.Schema{
		Type:       "object",
		Properties: oauth2Schema,
		Required:   []string{"type", "flows"},
	}
}
