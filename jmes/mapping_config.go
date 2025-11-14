package jmes

import (
	"encoding/json"
	"fmt"

	"github.com/hasura/goenvconf"
	"go.yaml.in/yaml/v4"
)

type FieldMappingConfigInterface interface {
	Type() FieldMappingType
	Evaluate() (FieldMapping, error)
}

type FieldMappingConfig struct {
	FieldMappingConfigInterface `yaml:",inline"`
}

type rawFieldMappingConfig struct {
	Type FieldMappingType `json:"type" yaml:"type"`
}

func NewFieldMappingConfig(inner FieldMappingConfigInterface) FieldMappingConfig {
	return FieldMappingConfig{FieldMappingConfigInterface: inner}
}

func (fm FieldMappingConfig) Interface() FieldMappingConfigInterface {
	return fm.FieldMappingConfigInterface
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

func (FieldMappingEntryConfig) Type() FieldMappingType {
	return FieldMappingTypeField
}

func (fm FieldMappingEntryConfig) IsZero() bool {
	return fm.Path == nil && fm.Default == nil
}

func (fm FieldMappingEntryConfig) Evaluate() (FieldMapping, error) {
	if fm.IsZero() {
		return FieldMapping{}, ErrFieldMappingEntryRequired
	}

	result := FieldMappingEntry{
		Path: fm.Path,
	}

	if fm.Default != nil {
		value, err := fm.Default.Get()
		if err != nil {
			return FieldMapping{}, err
		}

		result.Default = value
	}

	return NewFieldMapping(result), nil
}

// FieldMappingObjectConfig is a config for object mapping.
type FieldMappingObjectConfig struct {
	Properties map[string]FieldMappingConfig `json:"properties" yaml:"properties"`
}

func (FieldMappingObjectConfig) Type() FieldMappingType {
	return FieldMappingTypeObject
}

func (fm FieldMappingObjectConfig) IsZero() bool {
	return fm.Properties == nil
}

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
