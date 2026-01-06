package openapi

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/goutils"
	"gotest.tools/v3/assert"
)

func TestExtractCommonParametersOfOperation(t *testing.T) {
	testCases := []struct {
		name               string
		pathParams         []*highv3.Parameter
		operation          *highv3.Operation
		expectedPathParams []*highv3.Parameter
		expectedOpParams   []*highv3.Parameter
	}{
		{
			name:               "nil operation",
			pathParams:         []*highv3.Parameter{{Name: "id", In: string(InPath)}},
			operation:          nil,
			expectedPathParams: []*highv3.Parameter{{Name: "id", In: string(InPath)}},
			expectedOpParams:   nil,
		},
		{
			name:               "operation with no parameters",
			pathParams:         []*highv3.Parameter{{Name: "id", In: string(InPath)}},
			operation:          &highv3.Operation{Parameters: []*highv3.Parameter{}},
			expectedPathParams: []*highv3.Parameter{{Name: "id", In: string(InPath)}},
			expectedOpParams:   []*highv3.Parameter{},
		},
		{
			name:       "operation with duplicate path parameter",
			pathParams: []*highv3.Parameter{{Name: "id", In: string(InPath)}},
			operation: &highv3.Operation{
				Parameters: []*highv3.Parameter{
					{Name: "id", In: InPath},
					{Name: "filter", In: InQuery},
				},
			},
			expectedPathParams: []*highv3.Parameter{{Name: "id", In: InPath}},
			expectedOpParams:   []*highv3.Parameter{{Name: "filter", In: InQuery}},
		},
		{
			name:       "operation with new path parameter",
			pathParams: []*highv3.Parameter{{Name: "id", In: InPath}},
			operation: &highv3.Operation{
				Parameters: []*highv3.Parameter{
					{Name: "commentId", In: InPath},
					{Name: "filter", In: InQuery},
				},
			},
			expectedPathParams: []*highv3.Parameter{
				{Name: "id", In: InPath},
				{Name: "commentId", In: InPath},
			},
			expectedOpParams: []*highv3.Parameter{{Name: "filter", In: InQuery}},
		},
		{
			name:       "operation with query and header parameters",
			pathParams: []*highv3.Parameter{{Name: "id", In: InPath}},
			operation: &highv3.Operation{
				Parameters: []*highv3.Parameter{
					{Name: "filter", In: InQuery},
					{Name: "Authorization", In: InHeader},
				},
			},
			expectedPathParams: []*highv3.Parameter{{Name: "id", In: InPath}},
			expectedOpParams: []*highv3.Parameter{
				{Name: "filter", In: InQuery},
				{Name: "Authorization", In: InHeader},
			},
		},
		{
			name:       "operation with same name but different location",
			pathParams: []*highv3.Parameter{{Name: "id", In: InPath}},
			operation: &highv3.Operation{
				Parameters: []*highv3.Parameter{
					{Name: "id", In: InQuery},
				},
			},
			expectedPathParams: []*highv3.Parameter{{Name: "id", In: InPath}},
			expectedOpParams:   []*highv3.Parameter{{Name: "id", In: InQuery}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a copy of pathParams to avoid mutation affecting the test
			pathParamsCopy := make([]*highv3.Parameter, len(tc.pathParams))
			copy(pathParamsCopy, tc.pathParams)

			result := ExtractCommonParametersOfOperation(pathParamsCopy, tc.operation)

			assert.DeepEqual(t, tc.expectedPathParams, result, cmpopts.IgnoreUnexported(highv3.Parameter{}))
			if tc.operation != nil {
				assert.DeepEqual(t, tc.expectedOpParams, tc.operation.Parameters, cmpopts.IgnoreUnexported(highv3.Parameter{}))
			}
		})
	}
}

func TestMergeParameters(t *testing.T) {
	testCases := []struct {
		name     string
		dest     []*highv3.Parameter
		src      []*highv3.Parameter
		expected []*highv3.Parameter
	}{
		{
			name:     "empty dest and src",
			dest:     []*highv3.Parameter{},
			src:      []*highv3.Parameter{},
			expected: []*highv3.Parameter{},
		},
		{
			name:     "empty src",
			dest:     []*highv3.Parameter{{Name: "id", In: InPath}},
			src:      []*highv3.Parameter{},
			expected: []*highv3.Parameter{{Name: "id", In: InPath}},
		},
		{
			name: "empty dest",
			dest: []*highv3.Parameter{},
			src:  []*highv3.Parameter{{Name: "id", In: InPath}},
			expected: []*highv3.Parameter{
				{Name: "id", In: InPath},
			},
		},
		{
			name: "merge without duplicates",
			dest: []*highv3.Parameter{
				{Name: "id", In: InPath},
			},
			src: []*highv3.Parameter{
				{Name: "filter", In: InQuery},
			},
			expected: []*highv3.Parameter{
				{Name: "id", In: InPath},
				{Name: "filter", In: InQuery},
			},
		},
		{
			name: "merge with duplicate - src overrides dest",
			dest: []*highv3.Parameter{
				{Name: "id", In: InPath, Required: goutils.ToPtr(true)},
			},
			src: []*highv3.Parameter{
				{Name: "id", In: InPath, Required: goutils.ToPtr(false)},
			},
			expected: []*highv3.Parameter{
				{Name: "id", In: InPath, Required: goutils.ToPtr(false)},
			},
		},
		{
			name: "merge with same name but different location",
			dest: []*highv3.Parameter{
				{Name: "id", In: InPath},
			},
			src: []*highv3.Parameter{
				{Name: "id", In: InQuery},
			},
			expected: []*highv3.Parameter{
				{Name: "id", In: InPath},
				{Name: "id", In: InQuery},
			},
		},
		{
			name: "merge multiple parameters",
			dest: []*highv3.Parameter{
				{Name: "id", In: InPath},
				{Name: "filter", In: InQuery},
			},
			src: []*highv3.Parameter{
				{Name: "filter", In: InQuery, Required: goutils.ToPtr(true)},
				{Name: "sort", In: InQuery},
				{Name: "Authorization", In: InHeader},
			},
			expected: []*highv3.Parameter{
				{Name: "id", In: InPath},
				{Name: "filter", In: InQuery, Required: goutils.ToPtr(true)},
				{Name: "sort", In: InQuery},
				{Name: "Authorization", In: InHeader},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MergeParameters(tc.dest, tc.src)
			assert.DeepEqual(t, tc.expected, result, cmpopts.IgnoreUnexported(highv3.Parameter{}))
		})
	}
}
