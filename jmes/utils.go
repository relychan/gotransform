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
