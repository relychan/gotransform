package jmes

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath"
	"github.com/relychan/goutils"
)

// FieldMappingType represents a field mapping type enum.
type FieldMappingType string

const (
	FieldMappingTypeField  FieldMappingType = "field"
	FieldMappingTypeObject FieldMappingType = "object"
)

var (
	ErrUnsupportedFieldMappingType = errors.New("unsupported field mapping type")
	ErrFieldMappingEntryMalformed  = errors.New("field mapping entry is malformed")
	ErrFieldMappingEntryRequired   = errors.New("field mapping entry must not be empty")
	ErrFieldMappingObjectRequired  = errors.New("field mapping object must not be null")
)

// FieldMappingInterface abstracts a field mapping interface.
type FieldMappingInterface interface {
	goutils.IsZeroer

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

// IsZero checks if the config is empty.
func (fm FieldMapping) IsZero() bool {
	return fm.FieldMappingInterface == nil || fm.FieldMappingInterface.IsZero()
}

// Equal checks if this instance equals the target value.
func (fm FieldMapping) Equal(target FieldMapping) bool {
	if fm.FieldMappingInterface == target.FieldMappingInterface {
		return true
	}

	if fm.FieldMappingInterface == nil || target.FieldMappingInterface == nil {
		return false
	}

	if fm.Type() != target.Type() {
		return false
	}

	switch fmi := fm.FieldMappingInterface.(type) {
	case *FieldMappingEntry:
		return goutils.DeepEqual(fmi, target.FieldMappingInterface, true)
	case *FieldMappingEntryString:
		return goutils.DeepEqual(fmi, target.FieldMappingInterface, true)
	case *FieldMappingObject:
		return goutils.DeepEqual(fmi, target.FieldMappingInterface, true)
	default:
		return false
	}
}

// FieldMappingEntry is the entry to lookup field values with the specified JMES path.
type FieldMappingEntry struct {
	// Path is a JMESPath expression to find a value in the input data.
	Path *string
	// Default value to be used when no value is found when looking up the value using the path.
	Default any
}

var _ FieldMappingInterface = (*FieldMappingEntry)(nil)

// Type returns type of the field mapping entry.
func (FieldMappingEntry) Type() FieldMappingType {
	return FieldMappingTypeField
}

// IsZero checks if the config is empty.
func (fm FieldMappingEntry) IsZero() bool {
	return (fm.Path == nil || *fm.Path == "") && fm.Default == nil
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingEntry) Equal(target FieldMappingEntry) bool {
	return goutils.EqualComparablePtr(fm.Path, target.Path) &&
		goutils.DeepEqual(fm.Default, target.Default, false)
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

var _ FieldMappingInterface = (*FieldMappingObject)(nil)

// Type returns type of the field mapping entry.
func (FieldMappingObject) Type() FieldMappingType {
	return FieldMappingTypeObject
}

// IsZero checks if the config is empty.
func (fm FieldMappingObject) IsZero() bool {
	return len(fm.Properties) == 0
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingObject) Equal(target FieldMappingObject) bool {
	return goutils.EqualMap(fm.Properties, target.Properties, false)
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

var _ FieldMappingInterface = (*FieldMappingEntryString)(nil)

// Type returns type of the field mapping entry string.
func (FieldMappingEntryString) Type() FieldMappingType {
	return FieldMappingTypeField
}

// IsZero checks if the config is empty.
func (fm FieldMappingEntryString) IsZero() bool {
	return (fm.Path == nil || *fm.Path == "") && fm.Default == nil
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingEntryString) Equal(target FieldMappingEntryString) bool {
	return goutils.EqualComparablePtr(fm.Path, target.Path) &&
		goutils.EqualComparablePtr(fm.Default, target.Default)
}

// Evaluate validates and transforms data with the specified JMES path, returning any value.
func (fm FieldMappingEntryString) Evaluate(data any) (any, error) {
	return fm.EvaluateString(data)
}

// EvaluateString validates and transforms data with the specified JMES path, returning string value explicitly.
func (fm FieldMappingEntryString) EvaluateString(data any) (*string, error) {
	if fm.Path != nil {
		result, err := fm.evaluateStringFromPath(data)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if fm.Default != nil {
		return fm.Default, nil
	}

	return nil, nil
}

func (fm FieldMappingEntryString) evaluateStringFromPath(data any) (*string, error) {
	var result any

	if *fm.Path != "" {
		var err error

		result, err = jmespath.Search(*fm.Path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate mapping entry string: %w", err)
		}
	} else {
		result = data
	}

	if result == nil {
		return nil, nil
	}

	if str, ok := result.(string); ok {
		return &str, nil
	}

	return nil, fmt.Errorf(
		"%w, expected a string, got %s",
		ErrFieldMappingEntryMalformed,
		reflect.TypeOf(result),
	)
}
