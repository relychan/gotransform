// Package transformtypes defines common data types for other packages.
package transformtypes

import (
	"errors"
)

// TransformTemplateType represents the type of transform template enum.
type TransformTemplateType string

const (
	// TransformTemplateJMESPath is the transform template using JMESPath.
	TransformTemplateJMESPath TransformTemplateType = "jmespath"
	// TransformTemplateGo is the transform template using the standard text/template in Go.
	TransformTemplateGo TransformTemplateType = "gotmpl"
)

var (
	// ErrUnsupportedTransformerType occurs when the transformer type is not supported.
	ErrUnsupportedTransformerType = errors.New("unsupported transformer type")
	// ErrTemplateContentRequired occurs when the template content is empty.
	ErrTemplateContentRequired = errors.New("template content must not be empty")
)

// TemplateTransformerConfig abstracts the interface of a template transformer config.
type TemplateTransformerConfig interface {
	// Type returns type of the transformer.
	Type() TransformTemplateType
	// Validate checks if the config is valid.
	Validate() error
}
