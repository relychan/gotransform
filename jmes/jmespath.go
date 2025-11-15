// Package jmes implements the transform template using JMESPath templates.
package jmes

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

// Transform processes and injects data into the template to transform data.
func (jtt JMESTemplateTransformer) Transform(data any) (any, error) {
	return jtt.template.Evaluate(data)
}
