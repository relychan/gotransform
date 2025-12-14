package jmes

import (
	"encoding/json"

	"github.com/relychan/gotransform/transformtypes"
)

// JMESTransformerConfig represents configurations for the Go template transformer.
type JMESTransformerConfig struct {
	Template FieldMappingConfig `json:"template" yaml:"template"`
}

var _ transformtypes.TemplateTransformerConfig = (*JMESTransformerConfig)(nil)

// Type returns type of the transformer.
func (JMESTransformerConfig) Type() transformtypes.TransformTemplateType {
	return transformtypes.TransformTemplateJMESPath
}

// IsZero checks if the config is empty.
func (jt JMESTransformerConfig) IsZero() bool {
	return jt.Template.IsZero()
}

// Equal checks if this instance equals the target value.
func (jt JMESTransformerConfig) Equal(target JMESTransformerConfig) bool {
	return jt.Template.Equal(target.Template)
}

// Validate checks if the config is valid.
func (jt JMESTransformerConfig) Validate() error {
	if jt.Template.FieldMappingConfigInterface == nil {
		return transformtypes.ErrTemplateContentRequired
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (jt JMESTransformerConfig) MarshalJSON() ([]byte, error) {
	result := map[string]any{
		"type":     jt.Type(),
		"template": jt.Template,
	}

	return json.Marshal(result)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (jt JMESTransformerConfig) MarshalYAML() (any, error) {
	return map[string]any{
		"type":     jt.Type(),
		"template": jt.Template,
	}, nil
}
