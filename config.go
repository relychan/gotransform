package gotransform

import (
	"encoding/json"
	"fmt"

	"github.com/relychan/gotransform/gotmpl"
	"github.com/relychan/gotransform/jmes"
	"github.com/relychan/gotransform/transformtypes"
	"github.com/relychan/goutils"
	"go.yaml.in/yaml/v4"
)

// TemplateTransformerConfig represents configurations for transforming data.
type TemplateTransformerConfig struct {
	transformtypes.TemplateTransformerConfig `yaml:",inline"`
}

type rawTemplateTransformerConfig struct {
	Type transformtypes.TransformTemplateType `json:"type" yaml:"type"`
}

func (j TemplateTransformerConfig) Interface() transformtypes.TemplateTransformerConfig {
	return j.TemplateTransformerConfig
}

// IsZero checks if the config is empty.
func (j TemplateTransformerConfig) IsZero() bool {
	return j.TemplateTransformerConfig == nil || j.Type() == ""
}

// Equal checks if this instance equals the target value.
func (j TemplateTransformerConfig) Equal(target TemplateTransformerConfig) bool {
	if j.TemplateTransformerConfig == target.TemplateTransformerConfig {
		return true
	}

	if j.TemplateTransformerConfig == nil || target.TemplateTransformerConfig == nil {
		return false
	}

	switch c := j.TemplateTransformerConfig.(type) {
	case *gotmpl.GoTemplateTransformerConfig:
		return goutils.DeepEqual(c, target.TemplateTransformerConfig, true)
	case *jmes.JMESTransformerConfig:
		return goutils.DeepEqual(c, target.TemplateTransformerConfig, true)
	default:
		return false
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TemplateTransformerConfig) UnmarshalJSON(b []byte) error {
	var temp rawTemplateTransformerConfig

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	var config transformtypes.TemplateTransformerConfig

	switch temp.Type {
	case transformtypes.TransformTemplateGo:
		config = new(gotmpl.GoTemplateTransformerConfig)
	case transformtypes.TransformTemplateJMESPath:
		config = new(jmes.JMESTransformerConfig)
	default:
		return fmt.Errorf("%w: %s", transformtypes.ErrUnsupportedTransformerType, temp.Type)
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}

	j.TemplateTransformerConfig = config

	return nil
}

// UnmarshalYAML implements the custom behavior for the yaml.Unmarshaler interface.
func (j *TemplateTransformerConfig) UnmarshalYAML(value *yaml.Node) error {
	var temp rawTemplateTransformerConfig

	err := value.Decode(&temp)
	if err != nil {
		return err
	}

	var config transformtypes.TemplateTransformerConfig

	switch temp.Type {
	case transformtypes.TransformTemplateGo:
		config = new(gotmpl.GoTemplateTransformerConfig)
	case transformtypes.TransformTemplateJMESPath:
		config = new(jmes.JMESTransformerConfig)
	default:
		return fmt.Errorf("%w: %s", transformtypes.ErrUnsupportedTransformerType, temp.Type)
	}

	err = value.Decode(config)
	if err != nil {
		return err
	}

	j.TemplateTransformerConfig = config

	return nil
}
