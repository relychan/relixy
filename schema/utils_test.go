package schema

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestExtractCommonParametersOfOperation(t *testing.T) {
	testCases := []struct {
		name               string
		pathParams         []Parameter
		operation          *RelyProxyOperation
		expectedPathParams []Parameter
		expectedOpParams   []Parameter
	}{
		{
			name:               "nil operation",
			pathParams:         []Parameter{{Name: "id", In: InPath}},
			operation:          nil,
			expectedPathParams: []Parameter{{Name: "id", In: InPath}},
			expectedOpParams:   nil,
		},
		{
			name:               "operation with no parameters",
			pathParams:         []Parameter{{Name: "id", In: InPath}},
			operation:          &RelyProxyOperation{Parameters: []Parameter{}},
			expectedPathParams: []Parameter{{Name: "id", In: InPath}},
			expectedOpParams:   []Parameter{},
		},
		{
			name:       "operation with duplicate path parameter",
			pathParams: []Parameter{{Name: "id", In: InPath}},
			operation: &RelyProxyOperation{
				Parameters: []Parameter{
					{Name: "id", In: InPath},
					{Name: "filter", In: InQuery},
				},
			},
			expectedPathParams: []Parameter{{Name: "id", In: InPath}},
			expectedOpParams:   []Parameter{{Name: "filter", In: InQuery}},
		},
		{
			name:       "operation with new path parameter",
			pathParams: []Parameter{{Name: "id", In: InPath}},
			operation: &RelyProxyOperation{
				Parameters: []Parameter{
					{Name: "commentId", In: InPath},
					{Name: "filter", In: InQuery},
				},
			},
			expectedPathParams: []Parameter{
				{Name: "id", In: InPath},
				{Name: "commentId", In: InPath},
			},
			expectedOpParams: []Parameter{{Name: "filter", In: InQuery}},
		},
		{
			name:       "operation with query and header parameters",
			pathParams: []Parameter{{Name: "id", In: InPath}},
			operation: &RelyProxyOperation{
				Parameters: []Parameter{
					{Name: "filter", In: InQuery},
					{Name: "Authorization", In: InHeader},
				},
			},
			expectedPathParams: []Parameter{{Name: "id", In: InPath}},
			expectedOpParams: []Parameter{
				{Name: "filter", In: InQuery},
				{Name: "Authorization", In: InHeader},
			},
		},
		{
			name:       "operation with same name but different location",
			pathParams: []Parameter{{Name: "id", In: InPath}},
			operation: &RelyProxyOperation{
				Parameters: []Parameter{
					{Name: "id", In: InQuery},
				},
			},
			expectedPathParams: []Parameter{{Name: "id", In: InPath}},
			expectedOpParams:   []Parameter{{Name: "id", In: InQuery}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a copy of pathParams to avoid mutation affecting the test
			pathParamsCopy := make([]Parameter, len(tc.pathParams))
			copy(pathParamsCopy, tc.pathParams)

			result := ExtractCommonParametersOfOperation(pathParamsCopy, tc.operation)

			assert.DeepEqual(t, tc.expectedPathParams, result)
			if tc.operation != nil {
				assert.DeepEqual(t, tc.expectedOpParams, tc.operation.Parameters)
			}
		})
	}
}

func TestMergeParameters(t *testing.T) {
	testCases := []struct {
		name     string
		dest     []Parameter
		src      []Parameter
		expected []Parameter
	}{
		{
			name:     "empty dest and src",
			dest:     []Parameter{},
			src:      []Parameter{},
			expected: []Parameter{},
		},
		{
			name:     "empty src",
			dest:     []Parameter{{Name: "id", In: InPath}},
			src:      []Parameter{},
			expected: []Parameter{{Name: "id", In: InPath}},
		},
		{
			name: "empty dest",
			dest: []Parameter{},
			src:  []Parameter{{Name: "id", In: InPath}},
			expected: []Parameter{
				{Name: "id", In: InPath},
			},
		},
		{
			name: "merge without duplicates",
			dest: []Parameter{
				{Name: "id", In: InPath},
			},
			src: []Parameter{
				{Name: "filter", In: InQuery},
			},
			expected: []Parameter{
				{Name: "id", In: InPath},
				{Name: "filter", In: InQuery},
			},
		},
		{
			name: "merge with duplicate - src overrides dest",
			dest: []Parameter{
				{Name: "id", In: InPath, Required: boolPtr(true)},
			},
			src: []Parameter{
				{Name: "id", In: InPath, Required: boolPtr(false)},
			},
			expected: []Parameter{
				{Name: "id", In: InPath, Required: boolPtr(false)},
			},
		},
		{
			name: "merge with same name but different location",
			dest: []Parameter{
				{Name: "id", In: InPath},
			},
			src: []Parameter{
				{Name: "id", In: InQuery},
			},
			expected: []Parameter{
				{Name: "id", In: InPath},
				{Name: "id", In: InQuery},
			},
		},
		{
			name: "merge multiple parameters",
			dest: []Parameter{
				{Name: "id", In: InPath},
				{Name: "filter", In: InQuery},
			},
			src: []Parameter{
				{Name: "filter", In: InQuery, Required: boolPtr(true)},
				{Name: "sort", In: InQuery},
				{Name: "Authorization", In: InHeader},
			},
			expected: []Parameter{
				{Name: "id", In: InPath},
				{Name: "filter", In: InQuery, Required: boolPtr(true)},
				{Name: "sort", In: InQuery},
				{Name: "Authorization", In: InHeader},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MergeParameters(tc.dest, tc.src)
			assert.DeepEqual(t, tc.expected, result)
		})
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
