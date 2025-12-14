// Package jmes implements the transform template using JMESPath templates.
package jmes

import "github.com/relychan/gotransform/transformtypes"

// JMESTemplateTransformer implements the transform template using JMESPath templates.
type JMESTemplateTransformer struct {
	template FieldMapping
}

// NewJMESTemplateTransformer creates a new JMESTemplateTransformer instance.
func NewJMESTemplateTransformer(template FieldMapping) *JMESTemplateTransformer {
	return &JMESTemplateTransformer{
		template: template,
	}
}

// Type returns the transform template type of this instance.
func (JMESTemplateTransformer) Type() transformtypes.TransformTemplateType {
	return transformtypes.TransformTemplateJMESPath
}

// IsZero checks if the JMESTemplateTransformer is empty (zero-valued).
func (jtt JMESTemplateTransformer) IsZero() bool {
	return jtt.template.IsZero()
}

// Transform processes and injects data into the template to transform data.
func (jtt JMESTemplateTransformer) Transform(data any) (any, error) {
	return jtt.template.Evaluate(data)
}

// Equal checks if this instance equals the target value.
func (jtt JMESTemplateTransformer) Equal(target JMESTemplateTransformer) bool {
	return jtt.template.Equal(target.template)
}
