package graphqlhandler

import (
	"testing"

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
