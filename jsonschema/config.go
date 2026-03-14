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

package main

import (
	"github.com/invopop/jsonschema"
	"github.com/relychan/gotransform"
	"github.com/relychan/gotransform/transformtypes"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type TemplateTransformerConfig gotransform.TemplateTransformerConfig

// JSONSchema is used to generate a custom jsonschema.
func (TemplateTransformerConfig) JSONSchema() *jsonschema.Schema {
	jmesPathProps := orderedmap.New[string, *jsonschema.Schema]()
	jmesPathProps.Set("type", &jsonschema.Schema{
		Type:        "string",
		Description: "Template type to be used for transforming response",
		Enum:        []any{transformtypes.TransformTemplateJMESPath},
	})
	jmesPathProps.Set("template", &jsonschema.Schema{
		Description: "Template content to be transformed",
		Ref:         "#/$defs/FieldMappingConfig",
	})

	goTemplateProps := orderedmap.New[string, *jsonschema.Schema]()
	goTemplateProps.Set("type", &jsonschema.Schema{
		Type:        "string",
		Description: "Template type to be used for transforming response",
		Enum:        []any{transformtypes.TransformTemplateGo},
	})
	goTemplateProps.Set("contentType", &jsonschema.Schema{
		Description: "The expected content type to be transformed",
		Type:        "string",
	})
	goTemplateProps.Set("template", &jsonschema.Schema{
		Description: "Template content to be transformed",
		Type:        "string",
	})

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "object",
				Title:       "TemplateTransformerJMESPathConfig",
				Description: "Transform responses using the standard JMESPath template",
				Required:    []string{"type", "template"},
				Properties:  jmesPathProps,
			},
			{
				Type:        "object",
				Title:       "TemplateTransformerGoTemplateConfig",
				Description: "Transform responses using the standard Go template",
				Properties:  goTemplateProps,
				Required:    []string{"type", "template", "contentType"},
			},
		},
	}
}
