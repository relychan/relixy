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

package schema

import (
	"testing"

	"github.com/relychan/relixy/schema/baseschema"
	"github.com/stretchr/testify/assert"
)

// TestRelixyOpenAPIResource_GetMetadata tests the GetMetadata method
func TestRelixyOpenAPIResource_GetMetadata(t *testing.T) {
	resource := RelixyOpenAPIResource{
		BaseResourceModel: baseschema.BaseResourceModel{
			Metadata: baseschema.RelixyResourceMetadata{
				Name: "test-api",
			},
		},
	}

	metadata := resource.GetMetadata()
	assert.Equal(t, "test-api", metadata.Name)
}
