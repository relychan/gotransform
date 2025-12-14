package jmes

import (
	"encoding/json"
	"testing"

	"github.com/hasura/goenvconf"
	"go.yaml.in/yaml/v4"
)

func TestFieldMappingEntryConfig_Type(t *testing.T) {
	config := FieldMappingEntryConfig{}
	if config.Type() != FieldMappingTypeField {
		t.Errorf("expected type to be %s, got: %s", FieldMappingTypeField, config.Type())
	}
}

func TestFieldMappingEntryConfig_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		config := FieldMappingEntryConfig{}
		if !config.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("with path", func(t *testing.T) {
		path := "name"
		config := FieldMappingEntryConfig{Path: &path}
		if config.IsZero() {
			t.Error("expected IsZero to return false when path is set")
		}
	})

	t.Run("with default", func(t *testing.T) {
		defaultVal := goenvconf.NewEnvAny("", "test")
		config := FieldMappingEntryConfig{Default: &defaultVal}
		if config.IsZero() {
			t.Error("expected IsZero to return false when default is set")
		}
	})
}

func TestFieldMappingEntryConfig_Equal(t *testing.T) {
	t.Run("equal configs", func(t *testing.T) {
		path := "name"
		defaultVal := goenvconf.NewEnvAny("", "test")
		config1 := FieldMappingEntryConfig{Path: &path, Default: &defaultVal}
		config2 := FieldMappingEntryConfig{Path: &path, Default: &defaultVal}
		if !config1.Equal(config2) {
			t.Error("expected configs to be equal")
		}
	})

	t.Run("different paths", func(t *testing.T) {
		path1 := "name"
		path2 := "title"
		config1 := FieldMappingEntryConfig{Path: &path1}
		config2 := FieldMappingEntryConfig{Path: &path2}
		if config1.Equal(config2) {
			t.Error("expected configs to be different")
		}
	})
}

func TestFieldMappingEntryConfig_EvaluateEntry(t *testing.T) {
	t.Run("evaluate with path", func(t *testing.T) {
		path := "name"
		config := FieldMappingEntryConfig{Path: &path}

		entry, err := config.EvaluateEntryEnv()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if entry.Path == nil || *entry.Path != "name" {
			t.Errorf("expected path to be 'name', got: %v", entry.Path)
		}
	})

	t.Run("evaluate with default", func(t *testing.T) {
		path := "name"
		defaultVal := goenvconf.NewEnvAny("", "default")
		config := FieldMappingEntryConfig{Path: &path, Default: &defaultVal}

		entry, err := config.EvaluateEntryEnv()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if entry.Default != "default" {
			t.Errorf("expected default to be 'default', got: %v", entry.Default)
		}
	})

	t.Run("error with empty config", func(t *testing.T) {
		config := FieldMappingEntryConfig{}

		_, err := config.EvaluateEntryEnv()
		if err == nil {
			t.Fatal("expected error for empty config, got nil")
		}

		if err != ErrFieldMappingEntryRequired {
			t.Errorf("expected error to be ErrFieldMappingEntryRequired, got: %v", err)
		}
	})
}

func TestFieldMappingObjectConfig_Type(t *testing.T) {
	config := FieldMappingObjectConfig{}
	if config.Type() != FieldMappingTypeObject {
		t.Errorf("expected type to be %s, got: %s", FieldMappingTypeObject, config.Type())
	}
}

func TestFieldMappingObjectConfig_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		config := FieldMappingObjectConfig{}
		if !config.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("with properties", func(t *testing.T) {
		path := "name"
		config := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"field": NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &path}),
			},
		}
		if config.IsZero() {
			t.Error("expected IsZero to return false when properties are set")
		}
	})
}

func TestFieldMappingObjectConfig_Equal(t *testing.T) {
	t.Run("equal configs", func(t *testing.T) {
		path := "name"
		config1 := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"field": NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &path}),
			},
		}
		config2 := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"field": NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &path}),
			},
		}
		if !config1.Equal(config2) {
			t.Error("expected configs to be equal")
		}
	})

	t.Run("different properties", func(t *testing.T) {
		path1 := "name"
		path2 := "title"
		config1 := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"field": NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &path1}),
			},
		}
		config2 := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"field": NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &path2}),
			},
		}
		if config1.Equal(config2) {
			t.Error("expected configs to be different")
		}
	})
}

func TestFieldMappingObjectConfig_Evaluate(t *testing.T) {
	t.Run("evaluate simple object", func(t *testing.T) {
		namePath := "name"
		agePath := "age"
		config := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"userName": NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &namePath}),
				"userAge":  NewFieldMappingConfig(&FieldMappingEntryConfig{Path: &agePath}),
			},
		}

		mapping, err := config.EvaluateEnv()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if mapping.IsZero() {
			t.Error("expected mapping to be non-zero")
		}

		obj, ok := mapping.FieldMappingInterface.(FieldMappingObject)
		if !ok {
			t.Fatalf("expected mapping to be FieldMappingObject, got: %T", mapping.FieldMappingInterface)
		}

		if len(obj.Properties) != 2 {
			t.Errorf("expected 2 properties, got: %d", len(obj.Properties))
		}
	})

	t.Run("error with nil config", func(t *testing.T) {
		config := FieldMappingObjectConfig{}

		_, err := config.EvaluateEnv()
		if err == nil {
			t.Fatal("expected error for nil config, got nil")
		}

		if err != ErrFieldMappingObjectRequired {
			t.Errorf("expected error to be ErrFieldMappingObjectRequired, got: %v", err)
		}
	})

	t.Run("error with nil field mapping", func(t *testing.T) {
		config := FieldMappingObjectConfig{
			Properties: map[string]FieldMappingConfig{
				"field": {},
			},
		}

		_, err := config.EvaluateEnv()
		if err == nil {
			t.Fatal("expected error for nil field mapping, got nil")
		}
	})
}

func TestFieldMappingEntryStringConfig_Type(t *testing.T) {
	config := FieldMappingEntryStringConfig{}
	if config.Type() != FieldMappingTypeField {
		t.Errorf("expected type to be %s, got: %s", FieldMappingTypeField, config.Type())
	}
}

func TestFieldMappingEntryStringConfig_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		config := FieldMappingEntryStringConfig{}
		if !config.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("with path", func(t *testing.T) {
		path := "name"
		config := FieldMappingEntryStringConfig{Path: &path}
		if config.IsZero() {
			t.Error("expected IsZero to return false when path is set")
		}
	})
}

func TestFieldMappingEntryStringConfig_EvaluateString(t *testing.T) {
	t.Run("evaluate with path", func(t *testing.T) {
		path := "name"
		config := FieldMappingEntryStringConfig{Path: &path}

		entry, err := config.EvaluateString(goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if entry.Path == nil || *entry.Path != "name" {
			t.Errorf("expected path to be 'name', got: %v", entry.Path)
		}
	})

	t.Run("evaluate with default", func(t *testing.T) {
		path := "name"
		defaultVal := goenvconf.NewEnvString("", "default")
		config := FieldMappingEntryStringConfig{Path: &path, Default: &defaultVal}

		entry, err := config.EvaluateString(goenvconf.GetOSEnv)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if entry.Default == nil || *entry.Default != "default" {
			t.Errorf("expected default to be 'default', got: %v", entry.Default)
		}
	})

	t.Run("error with empty config", func(t *testing.T) {
		config := FieldMappingEntryStringConfig{}

		_, err := config.EvaluateString(goenvconf.GetOSEnv)
		if err == nil {
			t.Fatal("expected error for empty config, got nil")
		}

		if err != ErrFieldMappingEntryRequired {
			t.Errorf("expected error to be ErrFieldMappingEntryRequired, got: %v", err)
		}
	})
}

func TestFieldMappingConfig_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshal field type", func(t *testing.T) {
		jsonData := `{"type": "field", "path": "name"}`

		var config FieldMappingConfig
		err := json.Unmarshal([]byte(jsonData), &config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if config.Type() != FieldMappingTypeField {
			t.Errorf("expected type to be %s, got: %s", FieldMappingTypeField, config.Type())
		}

		entry, ok := config.FieldMappingConfigInterface.(*FieldMappingEntryConfig)
		if !ok {
			t.Fatalf("expected config to be FieldMappingEntryConfig, got: %T", config.FieldMappingConfigInterface)
		}

		if entry.Path == nil || *entry.Path != "name" {
			t.Errorf("expected path to be 'name', got: %v", entry.Path)
		}
	})

	t.Run("unmarshal object type", func(t *testing.T) {
		jsonData := `{"type": "object", "properties": {"field": {"type": "field", "path": "name"}}}`

		var config FieldMappingConfig
		err := json.Unmarshal([]byte(jsonData), &config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if config.Type() != FieldMappingTypeObject {
			t.Errorf("expected type to be %s, got: %s", FieldMappingTypeObject, config.Type())
		}
	})

	t.Run("error with unsupported type", func(t *testing.T) {
		jsonData := `{"type": "unsupported"}`

		var config FieldMappingConfig
		err := json.Unmarshal([]byte(jsonData), &config)
		if err == nil {
			t.Fatal("expected error for unsupported type, got nil")
		}
	})
}

func TestFieldMappingConfig_UnmarshalYAML(t *testing.T) {
	t.Run("unmarshal field type", func(t *testing.T) {
		yamlData := `
type: field
path: name
`

		var config FieldMappingConfig
		err := yaml.Unmarshal([]byte(yamlData), &config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if config.Type() != FieldMappingTypeField {
			t.Errorf("expected type to be %s, got: %s", FieldMappingTypeField, config.Type())
		}
	})

	t.Run("unmarshal object type", func(t *testing.T) {
		yamlData := `
type: object
properties:
  field:
    type: field
    path: name
`

		var config FieldMappingConfig
		err := yaml.Unmarshal([]byte(yamlData), &config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if config.Type() != FieldMappingTypeObject {
			t.Errorf("expected type to be %s, got: %s", FieldMappingTypeObject, config.Type())
		}
	})
}
