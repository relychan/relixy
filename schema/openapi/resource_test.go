package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/schema/base_schema"
	"gotest.tools/v3/assert"
)

// TestRelixyOpenAPIResource_GetMetadata tests the GetMetadata method
func TestRelixyOpenAPIResource_GetMetadata(t *testing.T) {
	resource := RelixyOpenAPIResource{
		BaseResourceModel: base_schema.BaseResourceModel{
			Metadata: base_schema.RelixyResourceMetadata{
				Name: "test-api",
			},
		},
	}

	metadata := resource.GetMetadata()
	assert.Equal(t, "test-api", metadata.Name)
}

// TestRelixyOpenAPIResourceDefinition_MarshalJSON tests JSON marshaling
func TestRelixyOpenAPIResourceDefinition_MarshalJSON(t *testing.T) {
	t.Run("with_spec_only", func(t *testing.T) {
		def := RelixyOpenAPIResourceDefinition{
			Spec: &highv3.Document{
				Info: &base.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
		}

		data, err := json.Marshal(def)
		assert.NilError(t, err)
		assert.Assert(t, len(data) > 0)

		// Verify it can be unmarshaled back
		var result map[string]any
		err = json.Unmarshal(data, &result)
		assert.NilError(t, err)
		assert.Assert(t, result["spec"] != nil)
	})

	t.Run("with_ref_only", func(t *testing.T) {
		def := RelixyOpenAPIResourceDefinition{
			Ref: "https://example.com/openapi.yaml",
		}

		data, err := json.Marshal(def)
		assert.NilError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		assert.NilError(t, err)
		assert.Equal(t, "https://example.com/openapi.yaml", result["ref"])
	})

	t.Run("with_ref_and_spec", func(t *testing.T) {
		def := RelixyOpenAPIResourceDefinition{
			Ref: "https://example.com/openapi.yaml",
			Spec: &highv3.Document{
				Info: &base.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
		}

		data, err := json.Marshal(def)
		assert.NilError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		assert.NilError(t, err)
		assert.Equal(t, "https://example.com/openapi.yaml", result["ref"])
		assert.Assert(t, result["spec"] != nil)
	})
}

// TestRelixyOpenAPIResourceDefinition_Build tests the Build method
func TestRelixyOpenAPIResourceDefinition_Build(t *testing.T) {
	ctx := context.Background()

	t.Run("spec_only_no_ref", func(t *testing.T) {
		def := RelixyOpenAPIResourceDefinition{
			Spec: &highv3.Document{
				Info: &base.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
		}

		doc, err := def.Build(ctx)
		assert.NilError(t, err)
		assert.Assert(t, doc != nil)
		assert.Equal(t, "Test API", doc.Info.Title)
	})

	t.Run("no_spec_no_ref_error", func(t *testing.T) {
		def := RelixyOpenAPIResourceDefinition{}

		doc, err := def.Build(ctx)
		assert.Assert(t, err != nil)
		assert.Equal(t, ErrResourceSpecRequired, err)
		assert.Assert(t, doc == nil)
	})

	t.Run("with_invalid_ref", func(t *testing.T) {
		def := RelixyOpenAPIResourceDefinition{
			Ref: "nonexistent/file.json",
		}

		doc, err := def.Build(ctx)
		assert.Assert(t, err != nil)
		assert.Assert(t, doc == nil)
	})

	t.Run("with_ref_swagger_v2", func(t *testing.T) {
		testCases := []struct {
			Ref string
		}{
			{
				Ref: "petstore2",
			},
		}

		for _, tc := range testCases {
			def := RelixyOpenAPIResourceDefinition{
				Ref: fmt.Sprintf("testdata/%s/swagger.json", tc.Ref),
			}

			doc, err := def.Build(ctx)
			assert.NilError(t, err)
			assert.Assert(t, doc != nil)
			assert.Assert(t, doc.Info != nil)

			rawYamlBytes, err := doc.Render()
			assert.NilError(t, err)

			expectedPath := fmt.Sprintf("testdata/%s/expected.yaml", tc.Ref)
			// assert.NilError(t, os.WriteFile(expectedPath, rawYamlBytes, 0664))

			newDoc, err := libopenapi.NewDocument(rawYamlBytes)
			assert.NilError(t, err)

			expectedBytes, err := os.ReadFile(expectedPath)
			assert.NilError(t, err)

			expectedRawDoc, err := libopenapi.NewDocument(expectedBytes)
			assert.NilError(t, err)

			changes, err := libopenapi.CompareDocuments(expectedRawDoc, newDoc)
			assert.NilError(t, err)
			assert.Assert(t, len(changes.GetAllChanges()) == 0)
		}
	})

	t.Run("with_ref_openapi_v3", func(t *testing.T) {
		testCases := []struct {
			Ref string
		}{
			{
				Ref: "petstore3",
			},
		}

		for _, tc := range testCases {
			def := RelixyOpenAPIResourceDefinition{
				Ref: fmt.Sprintf("testdata/%s/openapi.json", tc.Ref),
			}

			doc, err := def.Build(ctx)
			assert.NilError(t, err)
			assert.Assert(t, doc != nil)
			assert.Assert(t, doc.Info != nil)

			rawYamlBytes, err := doc.Render()
			assert.NilError(t, err)

			expectedPath := fmt.Sprintf("testdata/%s/expected.yaml", tc.Ref)
			// assert.NilError(t, os.WriteFile(expectedPath, rawYamlBytes, 0664))

			newDoc, err := libopenapi.NewDocument(rawYamlBytes)
			assert.NilError(t, err)

			expectedBytes, err := os.ReadFile(expectedPath)
			assert.NilError(t, err)

			expectedRawDoc, err := libopenapi.NewDocument(expectedBytes)
			assert.NilError(t, err)

			changes, err := libopenapi.CompareDocuments(expectedRawDoc, newDoc)
			assert.NilError(t, err)
			assert.Assert(t, len(changes.GetAllChanges()) == 0)
		}
	})
}
