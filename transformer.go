// Package gotransform implements the universal template transformer.
package gotransform

import (
	"fmt"

	"github.com/relychan/gotransform/gotmpl"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/gotransform/transformtypes"
)

type TemplateTransformer interface {
	// Transform processes and injects data into the template to transform data.
	Transform(data any) (any, error)
}

// NewTransformerFromConfig creates a template transformer from configuration.
func NewTransformerFromConfig(
	name string,
	config TemplateTransformerConfig,
) (TemplateTransformer, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	switch conf := config.Interface().(type) {
	case *jmes.JMESTransformerConfig:
		fieldMapping, err := conf.Template.Evaluate()
		if err != nil {
			return nil, err
		}

		return jmes.NewJMESTemplateTransformer(fieldMapping), nil
	case *gotmpl.GoTemplateTransformerConfig:
		return gotmpl.NewGoTemplateTransformer(name, conf)
	default:
		return nil, fmt.Errorf("%w: %s", transformtypes.ErrUnsupportedTransformerType, config.Type())
	}
}
