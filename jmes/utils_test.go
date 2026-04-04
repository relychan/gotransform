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
	"testing"

	"github.com/hasura/goenvconf"
)

func TestEvaluateObjectFieldMappingEntries(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		input := map[string]FieldMappingEntryConfig{}
		result, err := EvaluateObjectFieldMappingEntries(input, goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty result, got: %v", result)
		}
	})

	t.Run("valid entries", func(t *testing.T) {
		namePath := "name"
		agePath := "age"
		input := map[string]FieldMappingEntryConfig{
			"userName": {Path: &namePath},
			"userAge":  {Path: &agePath},
		}

		result, err := EvaluateObjectFieldMappingEntries(input, goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 entries, got: %d", len(result))
		}

		if result["userName"].Path == nil {
			t.Errorf("expected userName path to be 'name', got: %v", result["userName"].Path)
		}

		if result["userAge"].Path == nil {
			t.Errorf("expected userAge path to be 'age', got: %v", result["userAge"].Path)
		}
	})

	t.Run("entry with default value", func(t *testing.T) {
		namePath := "name"
		defaultValue := goenvconf.NewEnvAny("", "default")
		input := map[string]FieldMappingEntryConfig{
			"userName": {Path: &namePath, Default: &defaultValue},
		}

		result, err := EvaluateObjectFieldMappingEntries(input, goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 entry, got: %d", len(result))
		}

		if result["userName"].Default != "default" {
			t.Errorf("expected default to be 'default', got: %v", result["userName"].Default)
		}
	})

	t.Run("error with empty entry", func(t *testing.T) {
		input := map[string]FieldMappingEntryConfig{
			"userName": {},
		}

		_, err := EvaluateObjectFieldMappingEntries(input, goenvconf.GetOSEnv)
		if err == nil {
			t.Fatal("expected error for empty entry, got nil")
		}
	})
}

func TestEvaluateObjectFieldMappingStringEntries(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		input := map[string]FieldMappingEntryStringConfig{}
		result, err := EvaluateObjectFieldMappingStringEntries(input, goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty result, got: %v", result)
		}
	})

	t.Run("valid entries", func(t *testing.T) {
		namePath := "name"
		emailPath := "email"
		input := map[string]FieldMappingEntryStringConfig{
			"userName":  {Path: &namePath},
			"userEmail": {Path: &emailPath},
		}

		result, err := EvaluateObjectFieldMappingStringEntries(input, goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 entries, got: %d", len(result))
		}

		if result["userName"].Path == nil {
			t.Errorf("expected userName path to be 'name', got: %v", result["userName"].Path)
		}

		if result["userEmail"].Path == nil {
			t.Errorf("expected userEmail path to be 'email', got: %v", result["userEmail"].Path)
		}
	})

	t.Run("entry with default value", func(t *testing.T) {
		namePath := "name"
		defaultValue := goenvconf.NewEnvString("", "default")
		input := map[string]FieldMappingEntryStringConfig{
			"userName": {Path: &namePath, Default: &defaultValue},
		}

		result, err := EvaluateObjectFieldMappingStringEntries(input, goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 entry, got: %d", len(result))
		}

		if result["userName"].Default == nil || *result["userName"].Default != "default" {
			t.Errorf("expected default to be 'default', got: %v", result["userName"].Default)
		}
	})

	t.Run("error with empty entry", func(t *testing.T) {
		input := map[string]FieldMappingEntryStringConfig{
			"userName": {},
		}

		_, err := EvaluateObjectFieldMappingStringEntries(input, goenvconf.GetOSEnv)
		if err == nil {
			t.Fatal("expected error for empty entry, got nil")
		}
	})
}
