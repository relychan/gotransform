package jmes

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath"
)

// FieldMappingType represents a field mapping type enum.
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

// FieldMappingInterface abstracts a field mapping interface.
type FieldMappingInterface interface {
	Type() FieldMappingType
	Evaluate(data any) (any, error)
}

// FieldMapping is a wrapper of a field mapping interface to evaluate data.
type FieldMapping struct {
	FieldMappingInterface
}

// NewFieldMapping creates a new field mapping instance.
func NewFieldMapping(inner FieldMappingInterface) FieldMapping {
	return FieldMapping{FieldMappingInterface: inner}
}

// Interface returns the inner field mapping interface.
func (fm FieldMapping) Interface() FieldMappingInterface { //nolint:ireturn
	return fm.FieldMappingInterface
}

// FieldMappingEntry is the entry to lookup field values with the specified JMES path.
type FieldMappingEntry struct {
	// Path is a JMESPath expression to find a value in the input data.
	Path *string
	// Default value to be used when no value is found when looking up the value using the path.
	Default any
}

// Type returns type of the field mapping entry.
func (FieldMappingEntry) Type() FieldMappingType {
	return FieldMappingTypeField
}

// Evaluate validates and transforms data with the specified JMES path.
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

// Type returns type of the field mapping entry.
func (FieldMappingObject) Type() FieldMappingType {
	return FieldMappingTypeObject
}

// Evaluate validates and transforms data with the specified JMES path.
func (fm FieldMappingObject) Evaluate(data any) (any, error) {
	result := make(map[string]any)

	for key, field := range fm.Properties {
		if field.FieldMappingInterface == nil {
			return nil, nil
		}

		value, err := field.Evaluate(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}

		result[key] = value
	}

	return result, nil
}

// FieldMappingEntryString is the entry to lookup string values with the specified JMES path.
type FieldMappingEntryString struct {
	// Path is a JMESPath expression to find a value in the input data.
	Path *string
	// Default value to be used when no value is found when looking up the value using the path.
	Default *string
}


// Type returns type of the field mapping entry string.
func (FieldMappingEntryString) Type() FieldMappingType {
	return FieldMappingTypeField
}

// Evaluate validates and transforms data with the specified JMES path, returning a string value.
func (fm FieldMappingEntryString) Evaluate(data any) (any, error) {
	var result any
	if fm.Path != nil {
		if *fm.Path != "" {
			var err error
			result, err = jmespath.Search(*fm.Path, data)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate mapping entry string: %w", err)
			}
		} else {
			result = data
		}
		if result != nil {
			if str, ok := result.(string); ok {
				return str, nil
			}
		}
	}
	if fm.Default != nil {
		return *fm.Default, nil
	}
	return nil, nil
}
