// Package jmes implements the transform template using JMESPath templates.
package jmes

import (
	"fmt"
)

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
func (gtt JMESTemplateTransformer) Transform(data any) (any, error) {
	return transform(gtt.template, data)
}

func transform(template FieldMapping, data any) (any, error) {
	switch expr := template.Interface().(type) {
	case *FieldMappingEntry:
		return expr.Evaluate(data)
	case *FieldMappingObject:
		result := map[string]any{}

		for key, field := range expr.Properties {
			fieldValue, err := field.Evaluate(data)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}

			result[key] = fieldValue
		}

		return result, nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedFieldMappingType, expr)
	}
}
