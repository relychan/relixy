// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package baseschema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

func TestRelixyResourceMetadata_JSONMarshal(t *testing.T) {
	metadata := RelixyResourceMetadata{
		Name:        "test-resource",
		Description: "Test description",
		Instruction: "Test instruction",
	}

	data, err := json.Marshal(metadata)
	assert.NoError(t, err)

	var result RelixyResourceMetadata
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, metadata.Name, result.Name)
	assert.Equal(t, metadata.Description, result.Description)
	assert.Equal(t, metadata.Instruction, result.Instruction)
}

func TestRelixyResourceMetadata_YAMLMarshal(t *testing.T) {
	metadata := RelixyResourceMetadata{
		Name:        "test-resource",
		Description: "Test description",
		Instruction: "Test instruction",
	}

	data, err := yaml.Marshal(metadata)
	assert.NoError(t, err)

	var result RelixyResourceMetadata
	err = yaml.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, metadata.Name, result.Name)
	assert.Equal(t, metadata.Description, result.Description)
	assert.Equal(t, metadata.Instruction, result.Instruction)
}

func TestRelixyResourceMetadata_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyResourceMetadata)
	}{
		{
			name: "complete metadata",
			jsonData: `{
				"name": "test-api",
				"description": "Test API description",
				"instruction": "Test instruction"
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, meta *RelixyResourceMetadata) {
				assert.Equal(t, "test-api", meta.Name)
				assert.Equal(t, "Test API description", meta.Description)
				assert.Equal(t, "Test instruction", meta.Instruction)
			},
		},
		{
			name: "minimal metadata",
			jsonData: `{
				"name": "test-api"
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, meta *RelixyResourceMetadata) {
				assert.Equal(t, "test-api", meta.Name)
				assert.Equal(t, "", meta.Description)
				assert.Equal(t, "", meta.Instruction)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var meta RelixyResourceMetadata
			err := json.Unmarshal([]byte(tc.jsonData), &meta)

			if tc.expectError {
				assert.True(t, err != nil)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &meta)
				}
			}
		})
	}
}

func TestBaseResourceModel_GetBaseResource(t *testing.T) {
	model := BaseResourceModel{
		Version: "v1",
		Kind:    "OpenAPI",
		Metadata: RelixyResourceMetadata{
			Name:        "test-resource",
			Description: "Test description",
		},
	}

	result := model.GetBaseResource()
	assert.Equal(t, model.Version, result.Version)
	assert.Equal(t, model.Kind, result.Kind)
	assert.Equal(t, model.Metadata.Name, result.Metadata.Name)
	assert.Equal(t, model.Metadata.Description, result.Metadata.Description)
}

func TestBaseResourceModel_JSONMarshal(t *testing.T) {
	model := BaseResourceModel{
		Version: "v1",
		Kind:    "OpenAPI",
		Metadata: RelixyResourceMetadata{
			Name:        "test-resource",
			Description: "Test description",
		},
	}

	data, err := json.Marshal(model)
	assert.NoError(t, err)

	var result BaseResourceModel
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, model.Version, result.Version)
	assert.Equal(t, model.Kind, result.Kind)
	assert.Equal(t, model.Metadata.Name, result.Metadata.Name)
}

func TestBaseResourceModel_YAMLMarshal(t *testing.T) {
	model := BaseResourceModel{
		Version: "v1",
		Kind:    "OpenAPI",
		Metadata: RelixyResourceMetadata{
			Name:        "test-resource",
			Description: "Test description",
		},
	}

	data, err := yaml.Marshal(model)
	assert.NoError(t, err)

	var result BaseResourceModel
	err = yaml.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, model.Version, result.Version)
	assert.Equal(t, model.Kind, result.Kind)
	assert.Equal(t, model.Metadata.Name, result.Metadata.Name)
}

func TestBaseResourceModel_JSONUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *BaseResourceModel)
	}{
		{
			name: "complete resource model",
			jsonData: `{
				"version": "v1",
				"kind": "OpenAPI",
				"metadata": {
					"name": "test-api",
					"description": "Test description"
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, model *BaseResourceModel) {
				assert.Equal(t, "v1", model.Version)
				assert.Equal(t, OpenAPIKind, model.Kind)
				assert.Equal(t, "test-api", model.Metadata.Name)
				assert.Equal(t, "Test description", model.Metadata.Description)
			},
		},
		{
			name: "minimal resource model",
			jsonData: `{
				"version": "v1",
				"kind": "OpenAPI",
				"metadata": {
					"name": "test-api"
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, model *BaseResourceModel) {
				assert.Equal(t, "v1", model.Version)
				assert.Equal(t, OpenAPIKind, model.Kind)
				assert.Equal(t, "test-api", model.Metadata.Name)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var model BaseResourceModel
			err := json.Unmarshal([]byte(tc.jsonData), &model)

			if tc.expectError {
				assert.True(t, err != nil)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &model)
				}
			}
		})
	}
}

func TestBaseResourceModel_YAMLUnmarshal(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkFunc   func(*testing.T, *BaseResourceModel)
	}{
		{
			name: "complete resource model",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api
  description: Test description`,
			expectError: false,
			checkFunc: func(t *testing.T, model *BaseResourceModel) {
				assert.Equal(t, "v1", model.Version)
				assert.Equal(t, OpenAPIKind, model.Kind)
				assert.Equal(t, "test-api", model.Metadata.Name)
				assert.Equal(t, "Test description", model.Metadata.Description)
			},
		},
		{
			name: "minimal resource model",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api`,
			expectError: false,
			checkFunc: func(t *testing.T, model *BaseResourceModel) {
				assert.Equal(t, "v1", model.Version)
				assert.Equal(t, OpenAPIKind, model.Kind)
				assert.Equal(t, "test-api", model.Metadata.Name)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var model BaseResourceModel
			err := yaml.Unmarshal([]byte(tc.yamlData), &model)

			if tc.expectError {
				assert.True(t, err != nil)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &model)
				}
			}
		})
	}
}
