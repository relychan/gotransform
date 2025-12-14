package gotmpl

import (
	"encoding/json"
	"testing"

	"github.com/relychan/gotransform/transformtypes"
	"go.yaml.in/yaml/v4"
)

func TestGoTemplateTransformerConfig_Type(t *testing.T) {
	config := GoTemplateTransformerConfig{}
	if config.Type() != transformtypes.TransformTemplateGo {
		t.Errorf("expected type to be %s, got: %s", transformtypes.TransformTemplateGo, config.Type())
	}
}

func TestGoTemplateTransformerConfig_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		config := GoTemplateTransformerConfig{}
		if !config.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("non-zero with template only", func(t *testing.T) {
		config := GoTemplateTransformerConfig{
			Template: "{{.name}}",
		}
		if config.IsZero() {
			t.Error("expected IsZero to return false when template is set")
		}
	})

	t.Run("non-zero with contentType only", func(t *testing.T) {
		config := GoTemplateTransformerConfig{
			ContentType: "text/plain",
		}
		if config.IsZero() {
			t.Error("expected IsZero to return false when contentType is set")
		}
	})

	t.Run("non-zero with both fields", func(t *testing.T) {
		config := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.name}}",
		}
		if config.IsZero() {
			t.Error("expected IsZero to return false when both fields are set")
		}
	})
}

func TestGoTemplateTransformerConfig_Equal(t *testing.T) {
	t.Run("equal configs", func(t *testing.T) {
		config1 := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.name}}",
		}
		config2 := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.name}}",
		}
		if !config1.Equal(config2) {
			t.Error("expected configs to be equal")
		}
	})

	t.Run("different content types", func(t *testing.T) {
		config1 := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.name}}",
		}
		config2 := GoTemplateTransformerConfig{
			ContentType: "application/json",
			Template:    "{{.name}}",
		}
		if config1.Equal(config2) {
			t.Error("expected configs to be different")
		}
	})

	t.Run("different templates", func(t *testing.T) {
		config1 := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.name}}",
		}
		config2 := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.title}}",
		}
		if config1.Equal(config2) {
			t.Error("expected configs to be different")
		}
	})
}

func TestGoTemplateTransformerConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "{{.name}}",
		}
		if err := config.Validate(); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("empty template", func(t *testing.T) {
		config := GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    "",
		}
		err := config.Validate()
		if err == nil {
			t.Fatal("expected error for empty template, got nil")
		}
		if err != transformtypes.ErrTemplateContentRequired {
			t.Errorf("expected error to be ErrTemplateContentRequired, got: %v", err)
		}
	})
}

func TestGoTemplateTransformerConfig_MarshalJSON(t *testing.T) {
	config := GoTemplateTransformerConfig{
		ContentType: "text/plain",
		Template:    "{{.name}}",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["type"] != "gotmpl" {
		t.Errorf("expected type to be 'gotmpl', got: %v", result["type"])
	}

	if result["contentType"] != "text/plain" {
		t.Errorf("expected contentType to be 'text/plain', got: %v", result["contentType"])
	}

	if result["template"] != "{{.name}}" {
		t.Errorf("expected template to be '{{.name}}', got: %v", result["template"])
	}
}

func TestGoTemplateTransformerConfig_MarshalYAML(t *testing.T) {
	config := GoTemplateTransformerConfig{
		ContentType: "application/json",
		Template:    `{"message": "{{.text}}"}`,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	var result map[string]any
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["type"] != "gotmpl" {
		t.Errorf("expected type to be 'gotmpl', got: %v", result["type"])
	}

	if result["contentType"] != "application/json" {
		t.Errorf("expected contentType to be 'application/json', got: %v", result["contentType"])
	}

	if result["template"] != `{"message": "{{.text}}"}` {
		t.Errorf("expected template to match, got: %v", result["template"])
	}
}
