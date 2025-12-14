package jmes

import (
	"fmt"

	"github.com/hasura/goenvconf"
)

// EvaluateObjectFieldMappingEntries validate and resolve the entry mapping fields of an object.
func EvaluateObjectFieldMappingEntries(
	input map[string]FieldMappingEntryConfig,
	getEnvFunc goenvconf.GetEnvFunc,
) (map[string]FieldMappingEntry, error) {
	props := make(map[string]FieldMappingEntry)

	if len(input) == 0 {
		return props, nil
	}

	for key, envField := range input {
		field, err := envField.EvaluateEntry(getEnvFunc)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}

		props[key] = field
	}

	return props, nil
}

// EvaluateObjectFieldMappingStringEntries validate and resolve the entry mapping for string fields of an object.
func EvaluateObjectFieldMappingStringEntries(
	input map[string]FieldMappingEntryStringConfig,
	getEnvFunc goenvconf.GetEnvFunc,
) (map[string]FieldMappingEntryString, error) {
	props := make(map[string]FieldMappingEntryString)

	if len(input) == 0 {
		return props, nil
	}

	for key, envField := range input {
		field, err := envField.EvaluateString(getEnvFunc)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}

		props[key] = field
	}

	return props, nil
}
