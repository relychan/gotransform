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
