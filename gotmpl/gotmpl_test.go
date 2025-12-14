package gotmpl

import (
	"testing"
)

func TestNewGoTemplateTransformer(t *testing.T) {
	t.Run("success with text template", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "application/json",
			Template:    `{"message": "{{.name}}"}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if transformer == nil {
			t.Fatal("expected transformer to be non-nil")
		}

		if transformer.contentType != "application/json" {
			t.Errorf("expected contentType to be 'application/json', got: %s", transformer.contentType)
		}
	})

	t.Run("success with html template", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/html",
			Template:    `<h1>{{.title}}</h1>`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if transformer == nil {
			t.Fatal("expected transformer to be non-nil")
		}

		if transformer.contentType != "text/html" {
			t.Errorf("expected contentType to be 'text/html', got: %s", transformer.contentType)
		}
	})

	t.Run("success with sprig functions", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{.name | upper}}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if transformer == nil {
			t.Fatal("expected transformer to be non-nil")
		}
	})

	t.Run("error with invalid template", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{.name`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err == nil {
			t.Fatal("expected error for invalid template, got nil")
		}

		if transformer != nil {
			t.Errorf("expected transformer to be nil, got: %v", transformer)
		}
	})
}

func TestGoTemplateTransformer_Transform(t *testing.T) {
	t.Run("transform with JSON output", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "application/json",
			Template:    `{"message": "{{.name}}"}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		data := map[string]any{"name": "John"}
		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := map[string]any{"message": "John"}
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("expected result to be map[string]any, got: %T", result)
		}

		if resultMap["message"] != expected["message"] {
			t.Errorf("expected message to be %v, got: %v", expected["message"], resultMap["message"])
		}
	})

	t.Run("transform with text output", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `Hello {{.name}}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		data := map[string]any{"name": "World"}
		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := "Hello World"
		if result != expected {
			t.Errorf("expected result to be %q, got: %q", expected, result)
		}
	})

	t.Run("transform with HTML output", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/html",
			Template:    `<h1>{{.title}}</h1>`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		data := map[string]any{"title": "Welcome"}
		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := "<h1>Welcome</h1>"
		if result != expected {
			t.Errorf("expected result to be %q, got: %q", expected, result)
		}
	})

	t.Run("transform with sprig functions", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{.name | upper}}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		data := map[string]any{"name": "hello"}
		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := "HELLO"
		if result != expected {
			t.Errorf("expected result to be %q, got: %q", expected, result)
		}
	})

	t.Run("error with invalid JSON output", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "application/json",
			Template:    `invalid json`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		data := map[string]any{}
		_, err = transformer.Transform(data)
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})

	t.Run("error with template execution failure", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{fail "intentional error"}}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		data := map[string]any{}
		_, err = transformer.Transform(data)
		if err == nil {
			t.Fatal("expected error for template execution failure, got nil")
		}
	})
}

func TestGoTemplateTransformer_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		transformer := GoTemplateTransformer{}
		if !transformer.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("non-zero value", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{.name}}`,
		}

		transformer, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer: %v", err)
		}

		if transformer.IsZero() {
			t.Error("expected IsZero to return false for non-zero value")
		}
	})
}

func TestGoTemplateTransformer_Equal(t *testing.T) {
	t.Run("equal transformers", func(t *testing.T) {
		config := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{.name}}`,
		}

		transformer1, err := NewGoTemplateTransformer("test", config)
		if err != nil {
			t.Fatalf("failed to create transformer1: %v", err)
		}

		// Test equality with itself
		if !transformer1.Equal(*transformer1) {
			t.Error("expected transformer to be equal to itself")
		}
	})

	t.Run("different content types", func(t *testing.T) {
		config1 := &GoTemplateTransformerConfig{
			ContentType: "text/plain",
			Template:    `{{.name}}`,
		}

		config2 := &GoTemplateTransformerConfig{
			ContentType: "application/json",
			Template:    `{{.name}}`,
		}

		transformer1, err := NewGoTemplateTransformer("test", config1)
		if err != nil {
			t.Fatalf("failed to create transformer1: %v", err)
		}

		transformer2, err := NewGoTemplateTransformer("test", config2)
		if err != nil {
			t.Fatalf("failed to create transformer2: %v", err)
		}

		if transformer1.Equal(*transformer2) {
			t.Error("expected transformers to be different")
		}
	})
}
