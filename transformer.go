// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gotransform implements the universal template transformer.
package gotransform

import (
	"fmt"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform/gotmpl"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/gotransform/transformtypes"
	"github.com/relychan/goutils"
)

// TemplateTransformer abstracts the interface to transform data.
type TemplateTransformer interface {
	goutils.IsZeroer
	// Type returns the transform template type of this instance.
	Type() transformtypes.TransformTemplateType
	// Transform processes and injects data into the template to transform data.
	Transform(data any) (any, error)
}

// NewTransformerFromConfig creates a template transformer from configuration.
func NewTransformerFromConfig(
	name string,
	config TemplateTransformerConfig,
	getEnvFunc goenvconf.GetEnvFunc,
) (TemplateTransformer, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	switch conf := config.Interface().(type) {
	case *jmes.JMESTransformerConfig:
		fieldMapping, err := conf.Template.Evaluate(getEnvFunc)
		if err != nil {
			return nil, err
		}

		return jmes.NewJMESTemplateTransformer(fieldMapping), nil
	case *gotmpl.GoTemplateTransformerConfig:
		return gotmpl.NewGoTemplateTransformer(name, conf)
	default:
		return nil, fmt.Errorf(
			"%w: %s",
			transformtypes.ErrUnsupportedTransformerType,
			config.Type(),
		)
	}
}

// EqualTemplateTransformer checks if both template transformers are equal.
func EqualTemplateTransformer(a, b TemplateTransformer) bool {
	if a == b {
		return true
	}

	if a == nil || b == nil || a.Type() != b.Type() {
		return false
	}

	switch ta := a.(type) {
	case *jmes.JMESTemplateTransformer:
		return goutils.DeepEqual(*ta, b, true)
	case *gotmpl.GoTemplateTransformer:
		return goutils.DeepEqual(*ta, b, true)
	default:
		return false
	}
}
