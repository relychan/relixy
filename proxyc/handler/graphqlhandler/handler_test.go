package graphqlhandler

import (
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
						Expression: "param.name",
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
						Expression: "query.limit[0]",
					},
					"offset": {
						Expression: "query.offset[0]",
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
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Handler.resolveRequestVariables(&tc.TemplateData)
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.Expected, result)
		})
	}
}
