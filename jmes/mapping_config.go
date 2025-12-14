package jmes

import (
	"encoding/json"
	"fmt"

	"github.com/hasura/goenvconf"
	"github.com/relychan/goutils"
	"go.yaml.in/yaml/v4"
)

// FieldMappingConfigInterface abstracts the interface of a field mapping config.
type FieldMappingConfigInterface interface {
	Type() FieldMappingType
	Evaluate() (FieldMapping, error)
}

// FieldMappingConfig represents a generic field mapping config.
type FieldMappingConfig struct {
	FieldMappingConfigInterface `yaml:",inline"`
}

type rawFieldMappingConfig struct {
	Type FieldMappingType `json:"type" yaml:"type"`
}

// NewFieldMappingConfig creates a new FieldMappingConfig instance.
func NewFieldMappingConfig(inner FieldMappingConfigInterface) FieldMappingConfig {
	return FieldMappingConfig{FieldMappingConfigInterface: inner}
}

// Interface returns the underlying config interface.
func (fm FieldMappingConfig) Interface() FieldMappingConfigInterface {
	return fm.FieldMappingConfigInterface
}

// IsZero checks if the config is empty.
func (fm FieldMappingConfig) IsZero() bool {
	return fm.FieldMappingConfigInterface == nil
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingConfig) Equal(target FieldMappingConfig) bool {
	if fm.FieldMappingConfigInterface == target.FieldMappingConfigInterface {
		return true
	}

	if fm.FieldMappingConfigInterface == nil || target.FieldMappingConfigInterface == nil {
		return false
	}

	if fm.Type() != target.Type() {
		return false
	}

	switch fmi := fm.FieldMappingConfigInterface.(type) {
	case *FieldMappingEntryConfig:
		return goutils.DeepEqual(fmi, target.FieldMappingConfigInterface, true)
	case *FieldMappingEntryStringConfig:
		return goutils.DeepEqual(fmi, target.FieldMappingConfigInterface, true)
	case *FieldMappingObjectConfig:
		return goutils.DeepEqual(fmi, target.FieldMappingConfigInterface, true)
	default:
		return false
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (fm *FieldMappingConfig) UnmarshalJSON(b []byte) error {
	var temp rawFieldMappingConfig

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	var config FieldMappingConfigInterface

	switch temp.Type {
	case FieldMappingTypeObject:
		config = new(FieldMappingObjectConfig)
	case FieldMappingTypeField:
		config = new(FieldMappingEntryConfig)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFieldMappingType, temp.Type)
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}

	fm.FieldMappingConfigInterface = config

	return nil
}

// UnmarshalYAML implements the custom behavior for the yaml.Unmarshaler interface.
func (fm *FieldMappingConfig) UnmarshalYAML(value *yaml.Node) error {
	var temp rawFieldMappingConfig

	err := value.Decode(&temp)
	if err != nil {
		return err
	}

	var config FieldMappingConfigInterface

	switch temp.Type {
	case FieldMappingTypeObject:
		config = new(FieldMappingObjectConfig)
	case FieldMappingTypeField:
		config = new(FieldMappingEntryConfig)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFieldMappingType, temp.Type)
	}

	err = value.Decode(config)
	if err != nil {
		return err
	}

	fm.FieldMappingConfigInterface = config

	return nil
}

// FieldMappingEntryConfig is the entry config to lookup field values with the specified JMES path.
type FieldMappingEntryConfig struct {
	// Path is a JMESPath expression to find a value in the input data.
	Path *string `json:"path,omitempty" yaml:"path,omitempty"`
	// Default value to be used when no value is found when looking up the value using the path.
	Default *goenvconf.EnvAny `json:"default,omitempty" yaml:"default,omitempty"`
}

var _ FieldMappingConfigInterface = (*FieldMappingEntryConfig)(nil)

// Type returns the type of field mapping config.
func (FieldMappingEntryConfig) Type() FieldMappingType {
	return FieldMappingTypeField
}

// IsZero checks if the config is empty.
func (fm FieldMappingEntryConfig) IsZero() bool {
	return fm.Path == nil && fm.Default == nil
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingEntryConfig) Equal(target FieldMappingEntryConfig) bool {
	return goutils.EqualComparablePtr(fm.Path, target.Path) &&
		goutils.DeepEqual(fm.Default, target.Default, false)
}

// Evaluate converts the config to the field mapping instance.
func (fm FieldMappingEntryConfig) Evaluate() (FieldMapping, error) {
	entry, err := fm.EvaluateEntry()
	if err != nil {
		return FieldMapping{}, err
	}

	return NewFieldMapping(entry), nil
}

// EvaluateEntry converts the config to the field mapping entry instance.
func (fm FieldMappingEntryConfig) EvaluateEntry() (FieldMappingEntry, error) {
	if fm.IsZero() {
		return FieldMappingEntry{}, ErrFieldMappingEntryRequired
	}

	result := FieldMappingEntry{
		Path: fm.Path,
	}

	if fm.Default != nil {
		value, err := fm.Default.Get()
		if err != nil {
			return FieldMappingEntry{}, err
		}

		result.Default = value
	}

	return result, nil
}

// FieldMappingObjectConfig represents a config for object mapping.
type FieldMappingObjectConfig struct {
	Properties map[string]FieldMappingConfig `json:"properties" yaml:"properties"`
}

var _ FieldMappingConfigInterface = (*FieldMappingObjectConfig)(nil)

// Type returns the type of field mapping config.
func (FieldMappingObjectConfig) Type() FieldMappingType {
	return FieldMappingTypeObject
}

// IsZero checks if the config is empty.
func (fm FieldMappingObjectConfig) IsZero() bool {
	return fm.Properties == nil
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingObjectConfig) Equal(target FieldMappingObjectConfig) bool {
	return goutils.EqualMap(fm.Properties, target.Properties, true)
}

// Evaluate converts the config to the field mapping instance.
func (fm FieldMappingObjectConfig) Evaluate() (FieldMapping, error) {
	if fm.IsZero() {
		return FieldMapping{}, ErrFieldMappingObjectRequired
	}

	result := FieldMappingObject{
		Properties: make(map[string]FieldMapping),
	}

	for key, fieldConfig := range fm.Properties {
		if fieldConfig.FieldMappingConfigInterface == nil {
			return FieldMapping{}, fmt.Errorf("%s: %w", key, ErrFieldMappingEntryRequired)
		}

		field, err := fieldConfig.Evaluate()
		if err != nil {
			return FieldMapping{}, fmt.Errorf("%s: %w", key, err)
		}

		result.Properties[key] = field
	}

	return NewFieldMapping(result), nil
}

// FieldMappingEntryStringConfig is the entry config to lookup string values with the specified JMES path.
type FieldMappingEntryStringConfig struct {
	// Path is a JMESPath expression to find a value in the input data.
	Path *string `json:"path,omitempty" yaml:"path,omitempty"`
	// Default value to be used when no value is found when looking up the value using the path.
	Default *goenvconf.EnvString `json:"default,omitempty" yaml:"default,omitempty"`
}

var _ FieldMappingConfigInterface = (*FieldMappingEntryStringConfig)(nil)

// Type returns the type of field mapping config.
func (FieldMappingEntryStringConfig) Type() FieldMappingType {
	return FieldMappingTypeField
}

// IsZero checks if the config is empty.
func (fm FieldMappingEntryStringConfig) IsZero() bool {
	return fm.Path == nil && fm.Default == nil
}

// Equal checks if this instance equals the target value.
func (fm FieldMappingEntryStringConfig) Equal(target FieldMappingEntryStringConfig) bool {
	return goutils.EqualComparablePtr(fm.Path, target.Path) &&
		goutils.EqualPtr(fm.Default, target.Default)
}

// Evaluate converts the config to the field mapping instance.
func (fm FieldMappingEntryStringConfig) Evaluate() (FieldMapping, error) {
	inner, err := fm.EvaluateString()
	if err != nil {
		return FieldMapping{}, err
	}

	return NewFieldMapping(inner), nil
}

// EvaluateString converts the config to the field mapping instance.
func (fm FieldMappingEntryStringConfig) EvaluateString() (FieldMappingEntryString, error) {
	if fm.IsZero() {
		return FieldMappingEntryString{}, ErrFieldMappingEntryRequired
	}

	result := FieldMappingEntryString{
		Path: fm.Path,
	}

	if fm.Default != nil {
		value, err := fm.Default.Get()
		if err != nil {
			return FieldMappingEntryString{}, err
		}

		result.Default = &value
	}

	return result, nil
}
