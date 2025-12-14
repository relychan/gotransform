package gotmpl

import (
	"encoding/json"

	"github.com/relychan/gotransform/transformtypes"
)

// GoTemplateTransformerConfig represents configurations for the Go template transformer.
type GoTemplateTransformerConfig struct {
	ContentType string `json:"contentType" jsonschema:"default=application/json" yaml:"contentType"`
	Template    string `json:"template"    yaml:"template"`
}

var _ transformtypes.TemplateTransformerConfig = (*GoTemplateTransformerConfig)(nil)

// Type returns type of the transformer.
func (GoTemplateTransformerConfig) Type() transformtypes.TransformTemplateType {
	return transformtypes.TransformTemplateGo
}

// IsZero checks if the config is empty.
func (gt GoTemplateTransformerConfig) IsZero() bool {
	return gt.ContentType == "" && gt.Template == ""
}

// Equal checks if this instance equals the target value.
func (gt GoTemplateTransformerConfig) Equal(target GoTemplateTransformerConfig) bool {
	return gt.ContentType == target.ContentType &&
		gt.Template == target.Template
}

// Validate checks if the config is valid.
func (gt GoTemplateTransformerConfig) Validate() error {
	if gt.Template == "" {
		return transformtypes.ErrTemplateContentRequired
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (gt GoTemplateTransformerConfig) MarshalJSON() ([]byte, error) {
	result := map[string]any{
		"type":        gt.Type(),
		"contentType": gt.ContentType,
		"template":    gt.Template,
	}

	return json.Marshal(result)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (gt GoTemplateTransformerConfig) MarshalYAML() (any, error) {
	return map[string]any{
		"type":        gt.Type(),
		"contentType": gt.ContentType,
		"template":    gt.Template,
	}, nil
}
