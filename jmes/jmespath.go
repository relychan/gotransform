// Package jmes implements the transform template using JMESPath templates.
package jmes

import "github.com/jmespath-community/go-jmespath"

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

func transform(template any, data any) (any, error) {
	if template == nil {
		return nil, nil
	}

	switch expr := template.(type) {
	case string:
		if expr == "" {
			return "", nil
		}

		value, err := jmespath.Search(expr, data)
		if err != nil {
			return expr, err
		}

		return value, nil
	case map[string]any:
		result := make(map[string]any)

		for key, value := range expr {
			elemValue, _ := transform(value, data)
			result[key] = elemValue
		}

		return result, nil
	case []any:
		result := make([]any, len(expr))

		for i, value := range expr {
			elemValue, _ := transform(value, data)
			result[i] = elemValue
		}

		return result, nil
	default:
		return template, nil
	}
}
