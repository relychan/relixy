package base_schema

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRelixyActionType_Constants(t *testing.T) {
	assert.Equal(t, RelixyActionType("graphql"), ProxyTypeGraphQL)
	assert.Equal(t, RelixyActionType("rest"), ProxyTypeREST)
}

func TestRelixyActionType_JSONMarshal(t *testing.T) {
	testCases := []struct {
		name     string
		value    RelixyActionType
		expected string
	}{
		{
			name:     "graphql type",
			value:    ProxyTypeGraphQL,
			expected: `"graphql"`,
		},
		{
			name:     "rest type",
			value:    ProxyTypeREST,
			expected: `"rest"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			assert.NilError(t, err)
			assert.Equal(t, tc.expected, string(data))
		})
	}
}

func TestRelixyActionType_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected RelixyActionType
	}{
		{
			name:     "graphql type",
			jsonData: `"graphql"`,
			expected: ProxyTypeGraphQL,
		},
		{
			name:     "rest type",
			jsonData: `"rest"`,
			expected: ProxyTypeREST,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result RelixyActionType
			err := json.Unmarshal([]byte(tc.jsonData), &result)
			assert.NilError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPrimitiveType_Constants(t *testing.T) {
	assert.Equal(t, PrimitiveType("string"), String)
	assert.Equal(t, PrimitiveType("number"), Number)
	assert.Equal(t, PrimitiveType("integer"), Integer)
	assert.Equal(t, PrimitiveType("boolean"), Boolean)
	assert.Equal(t, PrimitiveType("array"), Array)
	assert.Equal(t, PrimitiveType("object"), Object)
	assert.Equal(t, PrimitiveType("null"), Null)
}

func TestPrimitiveType_JSONMarshal(t *testing.T) {
	testCases := []struct {
		name     string
		value    PrimitiveType
		expected string
	}{
		{
			name:     "string type",
			value:    String,
			expected: `"string"`,
		},
		{
			name:     "number type",
			value:    Number,
			expected: `"number"`,
		},
		{
			name:     "integer type",
			value:    Integer,
			expected: `"integer"`,
		},
		{
			name:     "boolean type",
			value:    Boolean,
			expected: `"boolean"`,
		},
		{
			name:     "array type",
			value:    Array,
			expected: `"array"`,
		},
		{
			name:     "object type",
			value:    Object,
			expected: `"object"`,
		},
		{
			name:     "null type",
			value:    Null,
			expected: `"null"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			assert.NilError(t, err)
			assert.Equal(t, tc.expected, string(data))
		})
	}
}

func TestPrimitiveType_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected PrimitiveType
	}{
		{
			name:     "string type",
			jsonData: `"string"`,
			expected: String,
		},
		{
			name:     "number type",
			jsonData: `"number"`,
			expected: Number,
		},
		{
			name:     "integer type",
			jsonData: `"integer"`,
			expected: Integer,
		},
		{
			name:     "boolean type",
			jsonData: `"boolean"`,
			expected: Boolean,
		},
		{
			name:     "array type",
			jsonData: `"array"`,
			expected: Array,
		},
		{
			name:     "object type",
			jsonData: `"object"`,
			expected: Object,
		},
		{
			name:     "null type",
			jsonData: `"null"`,
			expected: Null,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result PrimitiveType
			err := json.Unmarshal([]byte(tc.jsonData), &result)
			assert.NilError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPrimitiveType_JSONSchema(t *testing.T) {
	schema := String.JSONSchema()
	assert.Assert(t, schema != nil)
	assert.Equal(t, "string", schema.Type)
	assert.Equal(t, "A primitive type in OpenAPI specification", schema.Description)
	assert.Assert(t, len(schema.Enum) == 7)
}
