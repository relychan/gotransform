// Package main generates the jsonschema for the transformer config.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/relychan/gotransform/jmes"
)

func main() {
	err := jsonSchemaConfiguration()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for TemplateTransformerConfig: %w", err))
	}
}

func jsonSchemaConfiguration() error {
	r := new(jsonschema.Reflector)

	err := r.AddGoComments(
		"github.com/relychan/gotransform",
		"../gotransform",
		jsonschema.WithFullComment(),
	)
	if err != nil {
		return err
	}

	reflectSchema := r.Reflect(TemplateTransformerConfig{})

	for _, externalType := range []any{
		jmes.FieldMappingObjectConfig{},
		jmes.FieldMappingEntryConfig{},
	} {
		externalSchema := r.Reflect(externalType)

		for key, def := range externalSchema.Definitions {
			if _, ok := reflectSchema.Definitions[key]; !ok {
				reflectSchema.Definitions[key] = def
			}
		}
	}

	// custom schema types
	reflectSchema.Definitions["FieldMappingObjectConfig"].Properties.Set("type", &jsonschema.Schema{
		Description: "Type of the field mapping config",
		Type:        "string",
		Enum:        []any{jmes.FieldMappingTypeObject},
	})
	reflectSchema.Definitions["FieldMappingObjectConfig"].Required = append(
		reflectSchema.Definitions["FieldMappingObjectConfig"].Required,
		"type",
	)

	reflectSchema.Definitions["FieldMappingEntryConfig"].Properties.Set("type", &jsonschema.Schema{
		Description: "Type of the field mapping config",
		Type:        "string",
		Enum:        []any{jmes.FieldMappingTypeField},
	})
	reflectSchema.Definitions["FieldMappingEntryConfig"].Required = append(
		reflectSchema.Definitions["FieldMappingEntryConfig"].Required,
		"type",
	)

	reflectSchema.Definitions["FieldMappingConfig"] = &jsonschema.Schema{
		Description: "Represents a generic field mapping config",
		OneOf: []*jsonschema.Schema{
			{
				Description: "Mapping configurations for object fields",
				Ref:         "#/$defs/FieldMappingObjectConfig",
			},
			{
				Description: "The mapping configuration for an entry field",
				Ref:         "#/$defs/FieldMappingEntryConfig",
			},
		},
	}

	schemaBytes, err := json.MarshalIndent(reflectSchema, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile( //nolint:gosec
		"jsonschema/gotransform.schema.json",
		schemaBytes,
		0o644, //nolint:mnd
	)
}
