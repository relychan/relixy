package base_schema

import (
	"encoding/json"
	"testing"

	orderedmap "github.com/pb33f/ordered-map/v2"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestRelixyAction_JSONMarshal(t *testing.T) {
	testCases := []struct {
		name   string
		action RelixyAction
	}{
		{
			name: "REST action with path",
			action: RelixyAction{
				Type: ProxyTypeREST,
				Path: "/api/v1/users",
			},
		},
		{
			name: "GraphQL action",
			action: RelixyAction{
				Type: ProxyTypeGraphQL,
				Request: RelixyGraphQLRequestConfig{
					Query: "query { users { id name } }",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.action)
			assert.NilError(t, err)

			var result RelixyAction
			err = json.Unmarshal(data, &result)
			assert.NilError(t, err)
			assert.Equal(t, tc.action.Type, result.Type)
			assert.Equal(t, tc.action.Path, result.Path)
		})
	}
}

func TestRelixyAction_YAMLMarshal(t *testing.T) {
	action := RelixyAction{
		Type: ProxyTypeREST,
		Path: "/api/v1/users",
	}

	data, err := yaml.Marshal(action)
	assert.NilError(t, err)

	var result RelixyAction
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, action.Type, result.Type)
	assert.Equal(t, action.Path, result.Path)
}

func TestRelixyAction_JSONSchema(t *testing.T) {
	schema := RelixyAction{}.JSONSchema()
	assert.Assert(t, schema != nil)
	assert.Assert(t, len(schema.OneOf) == 2)

	// Check REST schema
	restSchema := schema.OneOf[0]
	assert.Equal(t, "object", restSchema.Type)
	assert.Equal(t, "RelixyActionREST", restSchema.Title)
	assert.Assert(t, len(restSchema.Required) == 1)
	assert.Equal(t, "type", restSchema.Required[0])

	// Check GraphQL schema
	graphqlSchema := schema.OneOf[1]
	assert.Equal(t, "object", graphqlSchema.Type)
	assert.Equal(t, "RelixyActionGraphQL", graphqlSchema.Title)
	assert.Assert(t, len(graphqlSchema.Required) == 2)
}

func TestGraphQLVariableDefinition_JSONMarshal(t *testing.T) {
	varDef := GraphQLVariableDefinition{
		Path: "$.user.id",
	}

	data, err := json.Marshal(varDef)
	assert.NilError(t, err)

	var result GraphQLVariableDefinition
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, varDef.Path, result.Path)
}

func TestRelixyGraphQLRequestConfig_JSONMarshal(t *testing.T) {
	variables := orderedmap.New[string, *GraphQLVariableDefinition]()
	variables.Set("userId", &GraphQLVariableDefinition{
		Path: "$.user.id",
	})

	config := RelixyGraphQLRequestConfig{
		Query:     "query($userId: ID!) { user(id: $userId) { id name } }",
		Variables: variables,
	}

	data, err := json.Marshal(config)
	assert.NilError(t, err)

	var result RelixyGraphQLRequestConfig
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, config.Query, result.Query)
	assert.Assert(t, result.Variables != nil)
}

func TestRelixyGraphQLRequestConfig_YAMLMarshal(t *testing.T) {
	config := RelixyGraphQLRequestConfig{
		Query: "query { users { id name } }",
	}

	data, err := yaml.Marshal(config)
	assert.NilError(t, err)

	var result RelixyGraphQLRequestConfig
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Equal(t, config.Query, result.Query)
}

func TestRelixyGraphQLRequestConfig_JSONSchema(t *testing.T) {
	schema := RelixyGraphQLRequestConfig{}.JSONSchema()
	assert.Assert(t, schema != nil)
	assert.Equal(t, "object", schema.Type)
	assert.Assert(t, len(schema.Required) == 1)
	assert.Equal(t, "query", schema.Required[0])
}

func TestRelixyGraphQLResponseConfig_IsZero(t *testing.T) {
	testCases := []struct {
		name     string
		config   RelixyGraphQLResponseConfig
		expected bool
	}{
		{
			name:     "empty config",
			config:   RelixyGraphQLResponseConfig{},
			expected: true,
		},
		{
			name: "config with http error code",
			config: RelixyGraphQLResponseConfig{
				HTTPErrorCode: func() *int { v := 400; return &v }(),
			},
			expected: false,
		},
		{
			name: "config with nil http error code",
			config: RelixyGraphQLResponseConfig{
				HTTPErrorCode: nil,
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.config.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRelixyGraphQLResponseConfig_JSONMarshal(t *testing.T) {
	errorCode := 400
	config := RelixyGraphQLResponseConfig{
		HTTPErrorCode: &errorCode,
	}

	data, err := json.Marshal(config)
	assert.NilError(t, err)

	var result RelixyGraphQLResponseConfig
	err = json.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Assert(t, result.HTTPErrorCode != nil)
	assert.Equal(t, 400, *result.HTTPErrorCode)
}

func TestRelixyGraphQLResponseConfig_YAMLMarshal(t *testing.T) {
	errorCode := 500
	config := RelixyGraphQLResponseConfig{
		HTTPErrorCode: &errorCode,
	}

	data, err := yaml.Marshal(config)
	assert.NilError(t, err)

	var result RelixyGraphQLResponseConfig
	err = yaml.Unmarshal(data, &result)
	assert.NilError(t, err)
	assert.Assert(t, result.HTTPErrorCode != nil)
	assert.Equal(t, 500, *result.HTTPErrorCode)
}
