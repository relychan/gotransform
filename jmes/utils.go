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
