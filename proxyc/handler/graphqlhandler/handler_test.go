package graphqlhandler

import (
	"net/url"
	"testing"

	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/vektah/gqlparser/ast"
	"gotest.tools/v3/assert"
)

func TestTransformRequest(t *testing.T) {
	testCases := []struct {
		Name         string
		Handler      GraphQLHandler
		TemplateData proxyhandler.RequestTemplateData
		Expected     map[string]any
	}{
		{
			Name:     "empty",
			Handler:  GraphQLHandler{},
			Expected: map[string]any{},
		},
		{
			Name: "param_simple",
			Handler: GraphQLHandler{
				variableDefinitions: ast.VariableDefinitionList{
					{
						Variable: "name",
					},
				},
				variables: map[string]jmes.FieldMappingEntry{
					"name": {
						Path: goutils.ToPtr("param.name"),
					},
				},
			},
			TemplateData: proxyhandler.RequestTemplateData{
				Params: map[string]string{
					"name": "Queen",
				},
			},
			Expected: map[string]any{
				"name": "Queen",
			},
		},
		{
			Name: "query_simple",
			Handler: GraphQLHandler{
				variableDefinitions: ast.VariableDefinitionList{
					{
						Variable: "limit",
					},
					{
						Variable: "offset",
					},
				},
				variables: map[string]jmes.FieldMappingEntry{
					"limit": {
						Path: goutils.ToPtr("query.limit[0]"),
					},
					"offset": {
						Path: goutils.ToPtr("query.offset[0]"),
					},
				},
			},
			TemplateData: proxyhandler.RequestTemplateData{
				QueryParams: map[string][]string{
					"limit":  {"10"},
					"offset": {"1"},
				},
			},
			Expected: map[string]any{
				"limit":  "10",
				"offset": "1",
			},
		},
		{
			Name: "with_default_value",
			Handler: GraphQLHandler{
				variableDefinitions: ast.VariableDefinitionList{
					{
						Variable: "status",
					},
				},
				variables: map[string]jmes.FieldMappingEntry{
					"status": {
						Path:    goutils.ToPtr("param.status"),
						Default: "active",
					},
				},
			},
			TemplateData: proxyhandler.RequestTemplateData{
				Params: map[string]string{},
			},
			Expected: map[string]any{
				"status": "active",
			},
		},
		{
			Name: "body_variable",
			Handler: GraphQLHandler{
				variableDefinitions: ast.VariableDefinitionList{
					{
						Variable: "body",
					},
				},
				variables: map[string]jmes.FieldMappingEntry{},
			},
			TemplateData: proxyhandler.RequestTemplateData{
				Body: map[string]any{
					"name": "test",
				},
			},
			Expected: map[string]any{
				"body": map[string]any{
					"name": "test",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Handler.resolveRequestVariables(&tc.TemplateData)
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.Expected, result)
		})
	}
}

func TestGraphQLHandler_Type(t *testing.T) {
	handler := &GraphQLHandler{}
	assert.Equal(t, "graphql", string(handler.Type()))
}

func TestRequestTemplateData_ToMap(t *testing.T) {
	data := proxyhandler.RequestTemplateData{
		Params: map[string]string{
			"id": "123",
		},
		QueryParams: url.Values{
			"limit": []string{"10"},
		},
		Headers: map[string]string{
			"authorization": "Bearer token",
		},
		Body: map[string]any{
			"name": "test",
		},
	}

	result := data.ToMap()

	assert.Assert(t, result["param"] != nil)
	assert.Assert(t, result["query"] != nil)
	assert.Assert(t, result["headers"] != nil)
	assert.Assert(t, result["body"] != nil)
}

func TestResolveRequestExtensions(t *testing.T) {
	testCases := []struct {
		name         string
		handler      GraphQLHandler
		templateData proxyhandler.RequestTemplateData
		expected     map[string]any
	}{
		{
			name: "empty extensions",
			handler: GraphQLHandler{
				extensions: map[string]jmes.FieldMappingEntry{},
			},
			templateData: proxyhandler.RequestTemplateData{},
			expected:     map[string]any{},
		},
		{
			name: "extension with path",
			handler: GraphQLHandler{
				extensions: map[string]jmes.FieldMappingEntry{
					"tracing": {
						Path: goutils.ToPtr("headers.x_trace_id"),
					},
				},
			},
			templateData: proxyhandler.RequestTemplateData{
				Headers: map[string]string{
					"x_trace_id": "trace-123",
				},
			},
			expected: map[string]any{
				"tracing": "trace-123",
			},
		},
		{
			name: "extension with default value",
			handler: GraphQLHandler{
				extensions: map[string]jmes.FieldMappingEntry{
					"version": {
						Default: "1.0",
					},
				},
			},
			templateData: proxyhandler.RequestTemplateData{},
			expected: map[string]any{
				"version": "1.0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.handler.resolveRequestExtensions(&tc.templateData)
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.expected, result)
		})
	}
}
