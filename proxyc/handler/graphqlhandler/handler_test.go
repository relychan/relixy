package graphqlhandler

import (
	"net/url"
	"testing"

	"github.com/vektah/gqlparser/ast"
	"gotest.tools/v3/assert"
)

func TestTransformRequest(t *testing.T) {
	testCases := []struct {
		Name         string
		Handler      GraphQLHandler
		TemplateData requestTemplateData
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
				variables: map[string]graphqlVariable{
					"name": {
						Path: "param.name",
					},
				},
			},
			TemplateData: requestTemplateData{
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
				variables: map[string]graphqlVariable{
					"limit": {
						Path: "query.limit[0]",
					},
					"offset": {
						Path: "query.offset[0]",
					},
				},
			},
			TemplateData: requestTemplateData{
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
				variables: map[string]graphqlVariable{
					"status": {
						Path:    "param.status",
						Default: "active",
					},
				},
			},
			TemplateData: requestTemplateData{
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
				variables: map[string]graphqlVariable{},
			},
			TemplateData: requestTemplateData{
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
	data := requestTemplateData{
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
		templateData requestTemplateData
		expected     map[string]any
	}{
		{
			name: "empty extensions",
			handler: GraphQLHandler{
				extensions: map[string]graphqlVariable{},
			},
			templateData: requestTemplateData{},
			expected:     map[string]any{},
		},
		{
			name: "extension with path",
			handler: GraphQLHandler{
				extensions: map[string]graphqlVariable{
					"tracing": {
						Path: "headers.x_trace_id",
					},
				},
			},
			templateData: requestTemplateData{
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
				extensions: map[string]graphqlVariable{
					"version": {
						Default: "1.0",
					},
				},
			},
			templateData: requestTemplateData{},
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
