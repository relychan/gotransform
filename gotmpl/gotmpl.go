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

// Package gotmpl implements the template transformer using Go template.
package gotmpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	htmltemplate "html/template"
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/relychan/gotransform/transformtypes"
)

const contentTypeHTML = "text/html"

// Template abstracts the interface for both text and html template implementation.
type Template interface {
	Execute(wr io.Writer, data any) error
}

// GoTemplateTransformer implements the template transformer using Go template.
type GoTemplateTransformer struct {
	contentType string
	template    Template
}

// NewGoTemplateTransformer creates a new GoTemplateTransformer instance.
func NewGoTemplateTransformer(
	name string,
	config *GoTemplateTransformerConfig,
) (*GoTemplateTransformer, error) {
	result := &GoTemplateTransformer{
		contentType: config.ContentType,
	}

	var err error

	if strings.HasPrefix(config.ContentType, contentTypeHTML) {
		result.template, err = htmltemplate.New(name).Funcs(sprig.FuncMap()).Parse(config.Template)
	} else {
		result.template, err = template.New(name).Funcs(sprig.FuncMap()).Parse(config.Template)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse template %q: %w", name, err)
	}

	return result, nil
}

// Type returns the transform template type of this instance.
func (GoTemplateTransformer) Type() transformtypes.TransformTemplateType {
	return transformtypes.TransformTemplateGo
}

// IsZero checks if the transformer is zero-valued.
func (gtt GoTemplateTransformer) IsZero() bool {
	return gtt.contentType == "" && gtt.template == nil
}

// Equal checks if this instance equals the target value.
func (gtt GoTemplateTransformer) Equal(target GoTemplateTransformer) bool {
	return gtt.contentType == target.contentType &&
		gtt.template == target.template
}

// Transform processes and injects data into the template to transform data.
func (gtt GoTemplateTransformer) Transform(data any) (any, error) {
	var buffer bytes.Buffer

	err := gtt.template.Execute(&buffer, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	switch gtt.contentType {
	case "application/json":
		var result any

		err := json.Unmarshal(buffer.Bytes(), &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON result: %w", err)
		}

		return result, nil
	default:
		return buffer.String(), nil
	}
}
