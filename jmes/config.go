package jmes

import "github.com/relychan/gotransform/transformtypes"

// JMESTransformerConfig represents configurations for the Go template transformer.
type JMESTransformerConfig struct {
	Template FieldMappingConfig `json:"template" yaml:"template"`
}

var _ transformtypes.TemplateTransformerConfig = (*JMESTransformerConfig)(nil)

// Type returns type of the transformer.
func (JMESTransformerConfig) Type() transformtypes.TransformTemplateType {
	return transformtypes.TransformTemplateJMESPath
}

// Validate checks if the config is valid.
func (gt JMESTransformerConfig) Validate() error {
	if gt.Template.FieldMappingConfigInterface == nil {
		return transformtypes.ErrTemplateContentRequired
	}

	return nil
}
