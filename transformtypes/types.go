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
