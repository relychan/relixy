package graphqlhandler

import (
	"testing"

	"github.com/hasura/goenvconf"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"gotest.tools/v3/assert"
)

func TestValidateGraphQLString(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		expectError   bool
		errorContains string
		checkHandler  func(t *testing.T, handler *GraphQLHandler)
	}{
		{
			name:          "empty query",
			query:         "",
			expectError:   true,
			errorContains: "query is required",
		},
		{
			name:          "invalid GraphQL syntax",
			query:         "query {",
			expectError:   true,
			errorContains: "Expected Name",
		},
		{
			name:        "valid simple query",
			query:       "query { users { id name } }",
			expectError: false,
			checkHandler: func(t *testing.T, handler *GraphQLHandler) {
				assert.Assert(t, handler != nil)
				assert.Equal(t, "query { users { id name } }", handler.query)
				assert.Equal(t, "query", string(handler.operation))
			},
		},
		{
			name:        "valid mutation",
			query:       "mutation CreateUser($name: String!) { createUser(name: $name) { id } }",
			expectError: false,
			checkHandler: func(t *testing.T, handler *GraphQLHandler) {
				assert.Assert(t, handler != nil)
				assert.Equal(t, "mutation", string(handler.operation))
				assert.Equal(t, 1, len(handler.variableDefinitions))
				assert.Equal(t, "name", handler.variableDefinitions[0].Variable)
			},
		},
		{
			name:        "query with operation name",
			query:       "query GetUsers { users { id } }",
			expectError: false,
			checkHandler: func(t *testing.T, handler *GraphQLHandler) {
				assert.Assert(t, handler != nil)
				assert.Equal(t, "GetUsers", handler.operationName)
			},
		},
		{
			name: "multiple operations (batch)",
			query: `
				query GetUsers { users { id } }
				query GetPosts { posts { id } }
			`,
			expectError:   true,
			errorContains: "batch is not supported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := ValidateGraphQLString(tc.query)

			if tc.expectError {
				assert.Assert(t, err != nil, "expected error but got nil")
				if tc.errorContains != "" {
					assert.ErrorContains(t, err, tc.errorContains)
				}
			} else {
				assert.NilError(t, err)
				if tc.checkHandler != nil {
					tc.checkHandler(t, handler)
				}
			}
		})
	}
}

func TestValidateGraphQLVariables(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      *orderedmap.OrderedMap[string, *GraphQLVariableDefinition]
		getEnvFunc  goenvconf.GetEnvFunc
		expectError bool
		checkResult func(t *testing.T, result map[string]graphqlVariable)
	}{
		{
			name:        "nil inputs",
			inputs:      nil,
			getEnvFunc:  goenvconf.GetOSEnv,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]graphqlVariable) {
				assert.Equal(t, 0, len(result))
			},
		},
		{
			name:        "empty inputs",
			inputs:      orderedmap.New[string, *GraphQLVariableDefinition](),
			getEnvFunc:  goenvconf.GetOSEnv,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]graphqlVariable) {
				assert.Equal(t, 0, len(result))
			},
		},
		{
			name: "variable with path only",
			inputs: func() *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] {
				m := orderedmap.New[string, *GraphQLVariableDefinition]()
				m.Set("userId", &GraphQLVariableDefinition{
					Path: "param.id",
				})
				return m
			}(),
			getEnvFunc:  goenvconf.GetOSEnv,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]graphqlVariable) {
				assert.Equal(t, 1, len(result))
				assert.Equal(t, "param.id", result["userId"].Path)
				assert.Assert(t, result["userId"].Default == nil)
			},
		},
		{
			name: "variable with default value",
			inputs: func() *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] {
				m := orderedmap.New[string, *GraphQLVariableDefinition]()
				defaultValue := goenvconf.NewEnvAnyValue("default-value")
				m.Set("status", &GraphQLVariableDefinition{
					Path:    "query.status",
					Default: &defaultValue,
				})
				return m
			}(),
			getEnvFunc:  goenvconf.GetOSEnv,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]graphqlVariable) {
				assert.Equal(t, 1, len(result))
				assert.Equal(t, "query.status", result["status"].Path)
				assert.Equal(t, "default-value", result["status"].Default)
			},
		},
		{
			name: "multiple variables",
			inputs: func() *orderedmap.OrderedMap[string, *GraphQLVariableDefinition] {
				m := orderedmap.New[string, *GraphQLVariableDefinition]()
				m.Set("id", &GraphQLVariableDefinition{
					Path: "param.id",
				})
				m.Set("name", &GraphQLVariableDefinition{
					Path: "body.name",
				})
				return m
			}(),
			getEnvFunc:  goenvconf.GetOSEnv,
			expectError: false,
			checkResult: func(t *testing.T, result map[string]graphqlVariable) {
				assert.Equal(t, 2, len(result))
				assert.Equal(t, "param.id", result["id"].Path)
				assert.Equal(t, "body.name", result["name"].Path)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validateGraphQLVariables(tc.inputs, tc.getEnvFunc)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				if tc.checkResult != nil {
					tc.checkResult(t, result)
				}
			}
		})
	}
}
