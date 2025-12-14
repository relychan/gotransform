package jmes

import (
	"testing"
)

func TestNewJMESTemplateTransformer(t *testing.T) {
	path := "name"
	template := NewFieldMapping(&FieldMappingEntry{Path: &path})

	transformer := NewJMESTemplateTransformer(template)
	if transformer == nil {
		t.Fatal("expected transformer to be non-nil")
	}

	if transformer.template.IsZero() {
		t.Error("expected template to be non-zero")
	}
}

func TestJMESTemplateTransformer_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		transformer := &JMESTemplateTransformer{}
		if !transformer.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("non-zero value", func(t *testing.T) {
		path := "name"
		template := NewFieldMapping(&FieldMappingEntry{Path: &path})
		transformer := NewJMESTemplateTransformer(template)

		if transformer.IsZero() {
			t.Error("expected IsZero to return false for non-zero value")
		}
	})
}

func TestJMESTemplateTransformer_Transform(t *testing.T) {
	t.Run("transform with simple path", func(t *testing.T) {
		path := "name"
		template := NewFieldMapping(&FieldMappingEntry{Path: &path})
		transformer := NewJMESTemplateTransformer(template)

		data := map[string]any{"name": "John"}
		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != "John" {
			t.Errorf("expected result to be 'John', got: %v", result)
		}
	})

	t.Run("transform with object mapping", func(t *testing.T) {
		namePath := "user.name"
		emailPath := "user.email"
		template := NewFieldMapping(&FieldMappingObject{
			Properties: map[string]FieldMapping{
				"userName":  NewFieldMapping(&FieldMappingEntry{Path: &namePath}),
				"userEmail": NewFieldMapping(&FieldMappingEntry{Path: &emailPath}),
			},
		})
		transformer := NewJMESTemplateTransformer(template)

		data := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"email": "jane@example.com",
			},
		}

		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("expected result to be map[string]any, got: %T", result)
		}

		if resultMap["userName"] != "Jane" {
			t.Errorf("expected userName to be 'Jane', got: %v", resultMap["userName"])
		}

		if resultMap["userEmail"] != "jane@example.com" {
			t.Errorf("expected userEmail to be 'jane@example.com', got: %v", resultMap["userEmail"])
		}
	})

	t.Run("transform with array projection", func(t *testing.T) {
		path := "users[*].name"
		template := NewFieldMapping(&FieldMappingEntry{Path: &path})
		transformer := NewJMESTemplateTransformer(template)

		data := map[string]any{
			"users": []map[string]any{
				{"name": "Alice"},
				{"name": "Bob"},
			},
		}

		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		resultSlice, ok := result.([]any)
		if !ok {
			t.Fatalf("expected result to be []any, got: %T", result)
		}

		if len(resultSlice) != 2 {
			t.Errorf("expected result length to be 2, got: %d", len(resultSlice))
		}

		if resultSlice[0] != "Alice" {
			t.Errorf("expected first element to be 'Alice', got: %v", resultSlice[0])
		}

		if resultSlice[1] != "Bob" {
			t.Errorf("expected second element to be 'Bob', got: %v", resultSlice[1])
		}
	})

	t.Run("transform with default value", func(t *testing.T) {
		path := "nonexistent"
		template := NewFieldMapping(&FieldMappingEntry{
			Path:    &path,
			Default: "default_value",
		})
		transformer := NewJMESTemplateTransformer(template)

		data := map[string]any{"name": "John"}
		result, err := transformer.Transform(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != "default_value" {
			t.Errorf("expected result to be 'default_value', got: %v", result)
		}
	})
}

func TestJMESTemplateTransformer_Equal(t *testing.T) {
	t.Run("equal transformers", func(t *testing.T) {
		path := "name"
		template := NewFieldMapping(&FieldMappingEntry{Path: &path})
		transformer1 := NewJMESTemplateTransformer(template)
		transformer2 := NewJMESTemplateTransformer(template)

		if !transformer1.Equal(*transformer2) {
			t.Error("expected transformers to be equal")
		}
	})

	t.Run("different transformers", func(t *testing.T) {
		path1 := "name"
		path2 := "title"
		template1 := NewFieldMapping(&FieldMappingEntry{Path: &path1})
		template2 := NewFieldMapping(&FieldMappingEntry{Path: &path2})
		transformer1 := NewJMESTemplateTransformer(template1)
		transformer2 := NewJMESTemplateTransformer(template2)

		if transformer1.Equal(*transformer2) {
			t.Error("expected transformers to be different")
		}
	})
}

