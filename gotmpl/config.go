package gotmpl

import "github.com/relychan/gotransform/transformtypes"

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

// Validate checks if the config is valid.
func (gt GoTemplateTransformerConfig) Validate() error {
	if gt.Template == "" {
		return transformtypes.ErrTemplateContentRequired
	}

	return nil
}
