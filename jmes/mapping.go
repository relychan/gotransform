package jmes

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath"
)

type FieldMappingType string

const (
	FieldMappingTypeField  FieldMappingType = "field"
	FieldMappingTypeObject FieldMappingType = "object"
)

var (
	ErrUnsupportedFieldMappingType = errors.New("unsupported field mapping type")
	ErrFieldMappingEntryRequired   = errors.New("field mapping entry must not be empty")
	ErrFieldMappingObjectRequired  = errors.New("field mapping object must not be null")
)

type FieldMappingInterface interface {
	Type() FieldMappingType
	Evaluate(data any) (any, error)
}

type FieldMapping struct {
	inner FieldMappingInterface
}

func NewFieldMapping(inner FieldMappingInterface) FieldMapping {
	return FieldMapping{inner: inner}
}

func (fm FieldMapping) Interface() FieldMappingInterface { //nolint:ireturn
	return fm.inner
}

// FieldMappingEntry is the entry to lookup field values with the specified JMES path.
type FieldMappingEntry struct {
	// JSON pointer to find the particular claim in the decoded JWT token.
	Path *string
	// Default value to be used when no value is found when looking up the value using the path.
	Default any
}

func (FieldMappingEntry) Type() FieldMappingType {
	return FieldMappingTypeField
}

func (fm FieldMappingEntry) Evaluate(data any) (any, error) {
	if fm.Path != nil {
		result := data

		if *fm.Path != "" {
			var err error

			result, err = jmespath.Search(*fm.Path, data)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate mapping entry: %w", err)
			}
		}

		if result != nil {
			return result, nil
		}
	}

	return fm.Default, nil
}

// FieldMappingObject is the entry to lookup object values with the specified JMES path.
type FieldMappingObject struct {
	Properties map[string]FieldMapping `json:"properties" yaml:"properties"`
}

func (FieldMappingObject) Type() FieldMappingType {
	return FieldMappingTypeObject
}

func (fm FieldMappingObject) Evaluate(data any) (any, error) {
	result := make(map[string]any)

	for key, field := range fm.Properties {
		if field.inner == nil {
			return nil, nil
		}

		value, err := field.inner.Evaluate(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}

		result[key] = value
	}

	return result, nil
}
