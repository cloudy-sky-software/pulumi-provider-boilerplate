// Copyright 2022-2024, Cloudy Sky Software LLC.
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

package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudy-sky-software/pulumi-xyz/provider/pkg/openapi"
)

func TestPulumiSchema(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "cmd", "pulumi-gen-xyz", "openapi.yml"))
	if err != nil {
		t.Fatalf("Failed reading openapi.yml: %v", err)
	}

	oaSpec := openapi.GetOpenAPISpec(b)

	PulumiSchema(*oaSpec)
}
