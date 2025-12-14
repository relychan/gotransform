package jmes

import (
	"encoding/json"
	"testing"

	"github.com/relychan/gotransform/transformtypes"
)

func TestJMESTransformerConfig_Type(t *testing.T) {
	config := JMESTransformerConfig{}
	if config.Type() != transformtypes.TransformTemplateJMESPath {
		t.Errorf("expected type to be %s, got: %s", transformtypes.TransformTemplateJMESPath, config.Type())
	}
}

func TestJMESTransformerConfig_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		config := JMESTransformerConfig{}
		if !config.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("non-zero value", func(t *testing.T) {
		path := "name"
		config := JMESTransformerConfig{
			Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
				Path: &path,
			}),
		}
		if config.IsZero() {
			t.Error("expected IsZero to return false for non-zero value")
		}
	})
}

func TestJMESTransformerConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		path := "name"
		config := JMESTransformerConfig{
			Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
				Path: &path,
			}),
		}
		if err := config.Validate(); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("empty template", func(t *testing.T) {
		config := JMESTransformerConfig{}
		err := config.Validate()
		if err == nil {
			t.Fatal("expected error for empty template, got nil")
		}
	})
}

func TestJMESTransformerConfig_MarshalJSON(t *testing.T) {
	path := "name"
	config := JMESTransformerConfig{
		Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
			Path: &path,
		}),
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["type"] != "jmespath" {
		t.Errorf("expected type to be 'jmespath', got: %v", result["type"])
	}

	if result["template"] == nil {
		t.Error("expected template to be non-nil")
	}
}

func TestJMESTransformerConfig_Equal(t *testing.T) {
	t.Run("equal configs", func(t *testing.T) {
		path := "name"
		config1 := JMESTransformerConfig{
			Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
				Path: &path,
			}),
		}
		config2 := JMESTransformerConfig{
			Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
				Path: &path,
			}),
		}
		if !config1.Equal(config2) {
			t.Error("expected configs to be equal")
		}
	})

	t.Run("different configs", func(t *testing.T) {
		path1 := "name"
		path2 := "title"
		config1 := JMESTransformerConfig{
			Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
				Path: &path1,
			}),
		}
		config2 := JMESTransformerConfig{
			Template: NewFieldMappingConfig(&FieldMappingEntryConfig{
				Path: &path2,
			}),
		}
		if config1.Equal(config2) {
			t.Error("expected configs to be different")
		}
	})
}
