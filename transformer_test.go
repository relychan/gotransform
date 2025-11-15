package gotransform

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestTransformerJSON(t *testing.T) {
	testCases := []struct {
		File     string
		Input    any
		Expected any
	}{
		{
			File: "testdata/jmes.json",
			Input: map[string]any{
				"authors": []map[string]any{
					{
						"name": "Anna",
					},
					{
						"name": "Tom",
					},
				},
			},
			Expected: map[string]any{
				"foo": "bar",
				"author": map[string]any{
					"names": []any{"Anna", "Tom"},
				},
			},
		},
		{
			File: "testdata/gotmpl.json",
			Input: map[string]any{
				"hello": "Hello world",
			},
			Expected: "<h1>Hello world</h1>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.File, func(t *testing.T) {
			rawBytes, err := os.ReadFile(tc.File)
			if err != nil {
				t.Fatalf("failed to read file: %s", err)
			}

			var config TemplateTransformerConfig

			err = json.Unmarshal(rawBytes, &config)
			if err != nil {
				t.Fatalf("failed to decode JSON: %s", err)
			}

			transformer, err := NewTransformerFromConfig("test", config)
			if err != nil {
				t.Fatal(err)
			}

			result, err := transformer.Transform(tc.Input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.Expected, result) {
				t.Fatalf("not equal, expected: %v, got: %v", tc.Expected, result)
			}
		})
	}
}

func TestTransformerYAML(t *testing.T) {
	testCases := []struct {
		File     string
		Input    any
		Expected any
	}{
		{
			File: "testdata/jmes.yaml",
			Input: map[string]any{
				"data": map[string]any{
					"authors": []string{"Jon", "Tony"},
				},
			},
			Expected: []string{"Jon", "Tony"},
		},
		{
			File: "testdata/gotmpl.yaml",
			Input: map[string]any{
				"data": map[string]any{
					"authors": []string{"Jon", "Tony"},
				},
			},
			Expected: map[string]any{
				"hello": "Jon",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.File, func(t *testing.T) {
			rawBytes, err := os.ReadFile(tc.File)
			if err != nil {
				t.Fatalf("failed to read file: %s", err)
			}

			var config TemplateTransformerConfig

			err = yaml.Unmarshal(rawBytes, &config)
			if err != nil {
				t.Fatalf("failed to decode YAML: %s", err)
			}

			transformer, err := NewTransformerFromConfig("test", config)
			if err != nil {
				t.Fatal(err)
			}

			result, err := transformer.Transform(tc.Input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.Expected, result) {
				t.Fatalf("not equal, expected: %v, got: %v", tc.Expected, result)
			}
		})
	}
}
