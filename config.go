package gotransform

import (
	"encoding/json"
	"fmt"

	"github.com/relychan/gotransform/gotmpl"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/gotransform/transformtypes"
	"go.yaml.in/yaml/v4"
)

// TemplateTransformerConfig represents configurations for transforming data.
type TemplateTransformerConfig struct {
	transformtypes.TemplateTransformerConfig `yaml:",inline"`
}

type rawTemplateTransformerConfig struct {
	Type transformtypes.TransformTemplateType `json:"type" yaml:"type"`
}

func (j TemplateTransformerConfig) Interface() transformtypes.TemplateTransformerConfig {
	return j.TemplateTransformerConfig
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TemplateTransformerConfig) UnmarshalJSON(b []byte) error {
	var temp rawTemplateTransformerConfig

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	var config transformtypes.TemplateTransformerConfig

	switch temp.Type {
	case transformtypes.TransformTemplateGo:
		config = new(gotmpl.GoTemplateTransformerConfig)
	case transformtypes.TransformTemplateJMESPath:
		config = new(jmes.JMESTransformerConfig)
	default:
		return fmt.Errorf("%w: %s", transformtypes.ErrUnsupportedTransformerType, temp.Type)
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}

	j.TemplateTransformerConfig = config

	return nil
}

// UnmarshalYAML implements the custom behavior for the yaml.Unmarshaler interface.
func (j *TemplateTransformerConfig) UnmarshalYAML(value *yaml.Node) error {
	var temp rawTemplateTransformerConfig

	err := value.Decode(&temp)
	if err != nil {
		return err
	}

	var config transformtypes.TemplateTransformerConfig

	switch temp.Type {
	case transformtypes.TransformTemplateGo:
		config = new(gotmpl.GoTemplateTransformerConfig)
	case transformtypes.TransformTemplateJMESPath:
		config = new(jmes.JMESTransformerConfig)
	default:
		return fmt.Errorf("%w: %s", transformtypes.ErrUnsupportedTransformerType, temp.Type)
	}

	err = value.Decode(config)
	if err != nil {
		return err
	}

	j.TemplateTransformerConfig = config

	return nil
}

// RelyTransformConfig represents configurations for transforming data.
// type RelyTransformConfig struct {
// 	Type        TransformTemplateType `json:"type" yaml:"type"`
// 	ContentType string                `json:"contentType,omitempty" yaml:"contentType,omitempty"`
// 	Template    json.RawMessage       `json:"template" yaml:"template"`
// }

// JSONSchema is used to generate a custom jsonschema.
// func (j RelyProxyTransformDataConfig) JSONSchema() *jsonschema.Schema {
// 	jmesPathProps := wk8orderedmap.New[string, *jsonschema.Schema]()
// 	jmesPathProps.Set("type", &jsonschema.Schema{
// 		Type:        "string",
// 		Description: "Template type to be used for transforming response",
// 		Enum:        []any{TransformTemplateJMESPath},
// 	})
// 	jmesPathProps.Set("body", &jsonschema.Schema{
// 		Description:          "Body template to be transformed",
// 		Type:                 "object",
// 		AdditionalProperties: jsonschema.TrueSchema,
// 	})

// 	goTemplateProps := wk8orderedmap.New[string, *jsonschema.Schema]()
// 	goTemplateProps.Set("type", &jsonschema.Schema{
// 		Type:        "string",
// 		Description: "Template type to be used for transforming response",
// 		Enum:        []any{TransformTemplateGo},
// 	})
// 	goTemplateProps.Set("template", &jsonschema.Schema{
// 		Description: "Body string template to be transformed",
// 		Type:        "string",
// 	})

// 	return &jsonschema.Schema{
// 		OneOf: []*jsonschema.Schema{
// 			{
// 				Type:        "object",
// 				Description: "Transform responses using the standard JMESPath template",
// 				Required:    []string{"type", "body"},
// 				Properties:  jmesPathProps,
// 			},
// 			{
// 				Type:        "object",
// 				Description: "Transform responses using the standard Go template",
// 				Properties:  goTemplateProps,
// 				Required:    []string{"type", "template"},
// 			},
// 		},
// 	}
// }
