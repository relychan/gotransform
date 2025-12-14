package jmes

import (
	"testing"

	"github.com/relychan/goutils"
)

func TestFieldMappingEntry_Type(t *testing.T) {
	entry := FieldMappingEntry{}
	if entry.Type() != FieldMappingTypeField {
		t.Errorf("expected type to be %s, got: %s", FieldMappingTypeField, entry.Type())
	}
}

func TestFieldMappingEntry_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		entry := FieldMappingEntry{}
		if !entry.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("with path", func(t *testing.T) {
		path := "data.name"
		entry := FieldMappingEntry{Path: &path}
		if entry.IsZero() {
			t.Error("expected IsZero to return false when path is set")
		}
	})

	t.Run("with empty path", func(t *testing.T) {
		path := ""
		entry := FieldMappingEntry{Path: &path}
		if !entry.IsZero() {
			t.Error("expected IsZero to return true when path is empty")
		}
	})

	t.Run("with default value", func(t *testing.T) {
		entry := FieldMappingEntry{Default: "default"}
		if entry.IsZero() {
			t.Error("expected IsZero to return false when default is set")
		}
	})
}

func TestFieldMappingEntry_Equal(t *testing.T) {
	t.Run("equal entries", func(t *testing.T) {
		path := "data.name"
		entry1 := FieldMappingEntry{Path: &path, Default: "test"}
		entry2 := FieldMappingEntry{Path: &path, Default: "test"}
		if !entry1.Equal(entry2) {
			t.Error("expected entries to be equal")
		}
	})

	t.Run("different paths", func(t *testing.T) {
		path1 := "data.name"
		path2 := "data.title"
		entry1 := FieldMappingEntry{Path: &path1}
		entry2 := FieldMappingEntry{Path: &path2}
		if entry1.Equal(entry2) {
			t.Error("expected entries to be different")
		}
	})

	t.Run("different defaults", func(t *testing.T) {
		path := "data.name"
		entry1 := FieldMappingEntry{Path: &path, Default: "test1"}
		entry2 := FieldMappingEntry{Path: &path, Default: "test2"}
		if entry1.Equal(entry2) {
			t.Error("expected entries to be different")
		}
	})
}

func TestFieldMappingEntry_Evaluate(t *testing.T) {
	t.Run("evaluate with path", func(t *testing.T) {
		path := "name"
		entry := FieldMappingEntry{Path: &path}
		data := map[string]any{"name": "John"}

		result, err := entry.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != "John" {
			t.Errorf("expected result to be 'John', got: %v", result)
		}
	})

	t.Run("evaluate with empty path", func(t *testing.T) {
		path := ""
		entry := FieldMappingEntry{Path: &path}
		data := map[string]any{"name": "John"}

		result, err := entry.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !goutils.DeepEqual(result, data, false) {
			t.Errorf("expected result to be the input data, got: %v", result)
		}
	})

	t.Run("evaluate with nested path", func(t *testing.T) {
		path := "user.name"
		entry := FieldMappingEntry{Path: &path}
		data := map[string]any{
			"user": map[string]any{
				"name": "Jane",
			},
		}

		result, err := entry.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != "Jane" {
			t.Errorf("expected result to be 'Jane', got: %v", result)
		}
	})

	t.Run("evaluate with array path", func(t *testing.T) {
		path := "users[*].name"
		entry := FieldMappingEntry{Path: &path}
		data := map[string]any{
			"users": []map[string]any{
				{"name": "Alice"},
				{"name": "Bob"},
			},
		}

		result, err := entry.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := []any{"Alice", "Bob"}
		if !goutils.DeepEqual(result, expected, false) {
			t.Errorf("expected result to be %v, got: %v", expected, result)
		}
	})

	t.Run("evaluate with default when path not found", func(t *testing.T) {
		path := "nonexistent"
		entry := FieldMappingEntry{Path: &path, Default: "default_value"}
		data := map[string]any{"name": "John"}

		result, err := entry.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != "default_value" {
			t.Errorf("expected result to be 'default_value', got: %v", result)
		}
	})

	t.Run("evaluate with default when path is nil", func(t *testing.T) {
		entry := FieldMappingEntry{Default: "default_value"}
		data := map[string]any{"name": "John"}

		result, err := entry.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != "default_value" {
			t.Errorf("expected result to be 'default_value', got: %v", result)
		}
	})

	t.Run("error with invalid path", func(t *testing.T) {
		path := "invalid[["
		entry := FieldMappingEntry{Path: &path}
		data := map[string]any{"name": "John"}

		_, err := entry.Evaluate(data)
		if err == nil {
			t.Fatal("expected error for invalid path, got nil")
		}
	})
}

func TestFieldMappingObject_Type(t *testing.T) {
	obj := FieldMappingObject{}
	if obj.Type() != FieldMappingTypeObject {
		t.Errorf("expected type to be %s, got: %s", FieldMappingTypeObject, obj.Type())
	}
}

func TestFieldMappingObject_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		obj := FieldMappingObject{}
		if !obj.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("with properties", func(t *testing.T) {
		path := "name"
		obj := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"field": NewFieldMapping(&FieldMappingEntry{Path: &path}),
			},
		}
		if obj.IsZero() {
			t.Error("expected IsZero to return false when properties are set")
		}
	})
}

func TestFieldMappingObject_Equal(t *testing.T) {
	t.Run("equal objects", func(t *testing.T) {
		path := "name"
		obj1 := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"field": NewFieldMapping(&FieldMappingEntry{Path: &path}),
			},
		}
		obj2 := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"field": NewFieldMapping(&FieldMappingEntry{Path: &path}),
			},
		}
		if !obj1.Equal(obj2) {
			t.Error("expected objects to be equal")
		}
	})

	t.Run("different properties", func(t *testing.T) {
		path1 := "name"
		path2 := "title"
		obj1 := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"field": NewFieldMapping(&FieldMappingEntry{Path: &path1}),
			},
		}
		obj2 := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"field": NewFieldMapping(&FieldMappingEntry{Path: &path2}),
			},
		}
		if obj1.Equal(obj2) {
			t.Error("expected objects to be different")
		}
	})
}

func TestFieldMappingObject_Evaluate(t *testing.T) {
	t.Run("evaluate simple object", func(t *testing.T) {
		namePath := "name"
		agePath := "age"
		obj := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"userName": NewFieldMapping(&FieldMappingEntry{Path: &namePath}),
				"userAge":  NewFieldMapping(&FieldMappingEntry{Path: &agePath}),
			},
		}
		data := map[string]any{
			"name": "John",
			"age":  30,
		}

		result, err := obj.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("expected result to be map[string]any, got: %T", result)
		}

		if resultMap["userName"] != "John" {
			t.Errorf("expected userName to be 'John', got: %v", resultMap["userName"])
		}

		if resultMap["userAge"] != 30 {
			t.Errorf("expected userAge to be 30, got: %v", resultMap["userAge"])
		}
	})

	t.Run("evaluate nested object", func(t *testing.T) {
		userNamePath := "user.name"
		userEmailPath := "user.email"
		innerObj := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"name":  NewFieldMapping(&FieldMappingEntry{Path: &userNamePath}),
				"email": NewFieldMapping(&FieldMappingEntry{Path: &userEmailPath}),
			},
		}

		data := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"email": "jane@example.com",
			},
		}

		result, err := innerObj.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("expected result to be map[string]any, got: %T", result)
		}

		if resultMap["name"] != "Jane" {
			t.Errorf("expected name to be 'Jane', got: %v", resultMap["name"])
		}

		if resultMap["email"] != "jane@example.com" {
			t.Errorf("expected email to be 'jane@example.com', got: %v", resultMap["email"])
		}
	})

	t.Run("error with nil field mapping", func(t *testing.T) {
		obj := FieldMappingObject{
			Properties: map[string]FieldMapping{
				"field": {},
			},
		}
		data := map[string]any{"name": "John"}

		result, err := obj.Evaluate(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result != nil {
			t.Errorf("expected result to be nil, got: %v", result)
		}
	})
}

func TestFieldMappingEntryString_Type(t *testing.T) {
	entry := FieldMappingEntryString{}
	if entry.Type() != FieldMappingTypeField {
		t.Errorf("expected type to be %s, got: %s", FieldMappingTypeField, entry.Type())
	}
}

func TestFieldMappingEntryString_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		entry := FieldMappingEntryString{}
		if !entry.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("with path", func(t *testing.T) {
		path := "data.name"
		entry := FieldMappingEntryString{Path: &path}
		if entry.IsZero() {
			t.Error("expected IsZero to return false when path is set")
		}
	})

	t.Run("with default", func(t *testing.T) {
		defaultVal := "default"
		entry := FieldMappingEntryString{Default: &defaultVal}
		if entry.IsZero() {
			t.Error("expected IsZero to return false when default is set")
		}
	})
}

func TestFieldMappingEntryString_Equal(t *testing.T) {
	t.Run("equal entries", func(t *testing.T) {
		path := "data.name"
		defaultVal := "test"
		entry1 := FieldMappingEntryString{Path: &path, Default: &defaultVal}
		entry2 := FieldMappingEntryString{Path: &path, Default: &defaultVal}
		if !entry1.Equal(entry2) {
			t.Error("expected entries to be equal")
		}
	})

	t.Run("different paths", func(t *testing.T) {
		path1 := "data.name"
		path2 := "data.title"
		entry1 := FieldMappingEntryString{Path: &path1}
		entry2 := FieldMappingEntryString{Path: &path2}
		if entry1.Equal(entry2) {
			t.Error("expected entries to be different")
		}
	})
}

func TestFieldMappingEntryString_EvaluateString(t *testing.T) {
	t.Run("evaluate with path", func(t *testing.T) {
		path := "name"
		entry := FieldMappingEntryString{Path: &path}
		data := map[string]any{"name": "John"}

		result, err := entry.EvaluateString(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result == nil || *result != "John" {
			t.Errorf("expected result to be 'John', got: %v", result)
		}
	})

	t.Run("evaluate with default", func(t *testing.T) {
		path := "nonexistent"
		defaultVal := "default"
		entry := FieldMappingEntryString{Path: &path, Default: &defaultVal}
		data := map[string]any{"name": "John"}

		result, err := entry.EvaluateString(data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result == nil || *result != "default" {
			t.Errorf("expected result to be 'default', got: %v", result)
		}
	})

	t.Run("error with non-string value", func(t *testing.T) {
		path := "age"
		entry := FieldMappingEntryString{Path: &path}
		data := map[string]any{"age": 30}

		_, err := entry.EvaluateString(data)
		if err == nil {
			t.Fatal("expected error for non-string value, got nil")
		}
	})
}

func TestFieldMapping_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		fm := FieldMapping{}
		if !fm.IsZero() {
			t.Error("expected IsZero to return true for zero value")
		}
	})

	t.Run("non-zero value", func(t *testing.T) {
		path := "name"
		fm := NewFieldMapping(&FieldMappingEntry{Path: &path})
		if fm.IsZero() {
			t.Error("expected IsZero to return false for non-zero value")
		}
	})
}

func TestFieldMapping_Equal(t *testing.T) {
	t.Run("equal field mappings", func(t *testing.T) {
		path := "name"
		fm1 := NewFieldMapping(&FieldMappingEntry{Path: &path})
		fm2 := NewFieldMapping(&FieldMappingEntry{Path: &path})
		if !fm1.Equal(fm2) {
			t.Error("expected field mappings to be equal")
		}
	})

	t.Run("different types", func(t *testing.T) {
		path := "name"
		fm1 := NewFieldMapping(&FieldMappingEntry{Path: &path})
		fm2 := NewFieldMapping(&FieldMappingObject{
			Properties: map[string]FieldMapping{},
		})
		if fm1.Equal(fm2) {
			t.Error("expected field mappings to be different")
		}
	})
}
