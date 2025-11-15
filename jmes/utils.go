package jmes

import "fmt"

// EvaluateObjectFieldMappingEntries validate and resolve the entry mapping fields of an object.
func EvaluateObjectFieldMappingEntries(
	input map[string]FieldMappingEntryConfig,
) (map[string]FieldMappingEntry, error) {
	props := make(map[string]FieldMappingEntry)

	if len(input) == 0 {
		return props, nil
	}

	for key, envField := range input {
		field, err := envField.EvaluateEntry()
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
) (map[string]FieldMappingEntryString, error) {
	props := make(map[string]FieldMappingEntryString)

	if len(input) == 0 {
		return props, nil
	}

	for key, envField := range input {
		field, err := envField.EvaluateString()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}

		props[key] = field
	}

	return props, nil
}
