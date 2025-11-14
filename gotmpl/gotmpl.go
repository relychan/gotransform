// Package gotmpl implements the template transformer using Go template.
package gotmpl

import (
	"bytes"
	htmltemplate "html/template"
	"io"
	"strings"
	"text/template"
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
		result.template, err = htmltemplate.New(name).Parse(config.Template)
	} else {
		result.template, err = template.New(name).Parse(config.Template)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Transform processes and injects data into the template to transform data.
func (gtt GoTemplateTransformer) Transform(data any) (any, error) {
	var buffer bytes.Buffer

	err := gtt.template.Execute(&buffer, data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
