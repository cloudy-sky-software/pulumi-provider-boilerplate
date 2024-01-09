// Copyright 2022-2024, Cloudy Sky Software LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openapi

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// GetOpenAPISpec returns a parsed and validated openapi doc.
func GetOpenAPISpec(data []byte) *openapi3.T {
	doc, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		contract.Failf("Failed to load openapi.yml: %v", err)
	}

	ctx := context.Background()
	// For the purposes of building a Pulumi schema, we don't care about
	// examples that may have been added to the spec by the cloud provider,
	// ignore those as those tend to have errors.
	if err := doc.Validate(ctx, openapi3.DisableExamplesValidation()); err != nil {
		contract.Failf("OpenAPI spec failed validation: %v", err)
	}

	return doc
}
