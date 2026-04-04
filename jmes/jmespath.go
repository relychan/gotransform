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
