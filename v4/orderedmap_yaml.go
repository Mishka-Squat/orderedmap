package orderedmap

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

type StringMapTree Of[string, any]

type yamlParser struct {
	V any
}

func YAML() yamlParser {
	return yamlParser{}
}

var _ yaml.Marshaler = (*yamlParser)(nil)
var _ yaml.Unmarshaler = (*yamlParser)(nil)

func (fs yamlParser) MarshalYAML() (interface{}, error) {
	return fs.V, nil
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

func (fs Of[K, V]) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	for k, v := range fs.All() {
		keyNode := &yaml.Node{}
		if err := keyNode.Encode(k); err != nil {
			return nil, err
		}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(v); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valueNode)
	}
	return node, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (fs *Of[K, V]) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.DocumentNode {
		return fs.UnmarshalYAML(value.Content[0])
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("orderedmap: cannot unmarshal yaml node of kind %v into Of[K, V]", value.Kind)
	}

	m := Make[K, V]()
	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valueNode := value.Content[i+1]

		var key K
		if err := keyNode.Decode(&key); err != nil {
			return err
		}

		var val V
		if err := valueNode.Decode(&val); err != nil {
			return err
		}

		m.Set(key, val)
	}
	*fs = m

	return nil
}
