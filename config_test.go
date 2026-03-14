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

package gotransform

import (
	"encoding/json"
	"testing"

	"github.com/relychan/gotransform/gotmpl"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/gotransform/transformtypes"
	"go.yaml.in/yaml/v4"
)

func TestTemplateTransformerConfig_UnmarshalJSON_GoTemplate(t *testing.T) {
	jsonData := `{
		"type": "gotmpl",
		"contentType": "text/plain",
		"template": "{{.name}}"
	}`

	var config TemplateTransformerConfig
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if config.Type() != transformtypes.TransformTemplateGo {
		t.Errorf("expected type to be %s, got: %s", transformtypes.TransformTemplateGo, config.Type())
	}

	goConfig, ok := config.Interface().(*gotmpl.GoTemplateTransformerConfig)
	if !ok {
		t.Fatalf("expected config to be GoTemplateTransformerConfig, got: %T", config.Interface())
	}

	if goConfig.ContentType != "text/plain" {
		t.Errorf("expected contentType to be 'text/plain', got: %s", goConfig.ContentType)
	}

	if goConfig.Template != "{{.name}}" {
		t.Errorf("expected template to be '{{.name}}', got: %s", goConfig.Template)
	}

	if config.IsZero() {
		t.Error("expected non-zero, got true")
	}

	if !config.Equal(TemplateTransformerConfig{
		TemplateTransformerConfig: goConfig,
	}) {
		t.Error("expected equal, got false")
	}

	if config.Equal(TemplateTransformerConfig{}) {
		t.Error("expected not-equal, got equal")
	}
}

func TestTemplateTransformerConfig_UnmarshalJSON_JMESPath(t *testing.T) {
	jsonData := `{
		"type": "jmespath",
		"template": {
			"type": "field",
			"path": "name"
		}
	}`

	var config TemplateTransformerConfig
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if config.Type() != transformtypes.TransformTemplateJMESPath {
		t.Errorf("expected type to be %s, got: %s", transformtypes.TransformTemplateJMESPath, config.Type())
	}

	jmesConfig, ok := config.Interface().(*jmes.JMESTransformerConfig)
	if !ok {
		t.Fatalf("expected config to be JMESTransformerConfig, got: %T", config.Interface())
	}

	if jmesConfig.Template.IsZero() {
		t.Error("expected template to be non-zero")
	}
}

func TestTemplateTransformerConfig_UnmarshalJSON_UnsupportedType(t *testing.T) {
	jsonData := `{
		"type": "unsupported"
	}`

	var config TemplateTransformerConfig
	err := json.Unmarshal([]byte(jsonData), &config)
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}

func TestTemplateTransformerConfig_UnmarshalYAML_GoTemplate(t *testing.T) {
	yamlData := `
type: gotmpl
contentType: application/json
template: '{"message": "{{.text}}"}'
`

	var config TemplateTransformerConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if config.Type() != transformtypes.TransformTemplateGo {
		t.Errorf("expected type to be %s, got: %s", transformtypes.TransformTemplateGo, config.Type())
	}

	goConfig, ok := config.Interface().(*gotmpl.GoTemplateTransformerConfig)
	if !ok {
		t.Fatalf("expected config to be GoTemplateTransformerConfig, got: %T", config.Interface())
	}

	if goConfig.ContentType != "application/json" {
		t.Errorf("expected contentType to be 'application/json', got: %s", goConfig.ContentType)
	}
}

func TestTemplateTransformerConfig_UnmarshalYAML_JMESPath(t *testing.T) {
	yamlData := `
type: jmespath
template:
  type: field
  path: name
`

	var config TemplateTransformerConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if config.Type() != transformtypes.TransformTemplateJMESPath {
		t.Errorf("expected type to be %s, got: %s", transformtypes.TransformTemplateJMESPath, config.Type())
	}

	jmesConfig, ok := config.Interface().(*jmes.JMESTransformerConfig)
	if !ok {
		t.Fatalf("expected config to be JMESTransformerConfig, got: %T", config.Interface())
	}

	if jmesConfig.Template.IsZero() {
		t.Error("expected template to be non-zero")
	}
}

func TestTemplateTransformerConfig_UnmarshalYAML_UnsupportedType(t *testing.T) {
	yamlData := `
type: unsupported
`

	var config TemplateTransformerConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}
