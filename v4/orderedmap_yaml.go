package orderedmap

import (
	"go.yaml.in/yaml/v4"
)

type StringMapTree Of[string, any]

type yamlParser struct {
	V any
}

func YAML() yamlParser {
	return yamlParser{}
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (fs *yamlParser) UnmarshalYAML(value *yaml.Node) error {
	// A YAML object/map is represented as a MappingNode
	switch value.Kind {
	case yaml.DocumentNode:
		return fs.UnmarshalYAML(value.Content[0])
	case yaml.SequenceNode:
		sequence := []any{}
		for i := 0; i < len(value.Content); i++ {
			valueNode := value.Content[i]

			var val any
			if valueNode.Kind == yaml.MappingNode || valueNode.Kind == yaml.SequenceNode {
				var yp yamlParser
				if err := valueNode.Decode(&yp); err != nil {
					return err
				}
				val = yp.V
			} else {
				if err := valueNode.Decode(&val); err != nil {
					return err
				}
			}
			sequence = append(sequence, val)
		}
		fs.V = sequence
	case yaml.MappingNode:
		mapping := Make[string, any]()
		// Content contains pairs of keys and values sequentially
		for i := 0; i < len(value.Content); i += 2 {
			keyNode := value.Content[i]
			valueNode := value.Content[i+1]

			var val any
			if valueNode.Kind == yaml.MappingNode || valueNode.Kind == yaml.SequenceNode {
				var yp yamlParser
				if err := valueNode.Decode(&yp); err != nil {
					return err
				}
				val = yp.V
			} else {
				if err := valueNode.Decode(&val); err != nil {
					return err
				}
			}
			mapping.Set(keyNode.Value, val)
		}
		fs.V = mapping
	}

	return nil
}
