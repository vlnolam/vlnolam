package config

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v4"
)

type Port uint16

type Ports struct {
	Items []PortSpec
}

func (p Ports) String() string {
	parts := make([]string, len(p.Items))
	for i, item := range p.Items {
		parts[i] = fmt.Sprintf("%+v", item) // Uses item's String() if exists
	}
	return strings.Join(parts, ", ")
}

func (p *Ports) UnmarshalYAML(node *yaml.Node) error {
	var spec PortSpec

	if node.Kind != yaml.SequenceNode {
		// Not an array -> single PortSpec
		err := unmarshalPortSpec(node, &spec)
		if err != nil {
			return err
		}

		p.Items = []PortSpec{spec}
		return nil
	}

	p.Items = []PortSpec{}

	// Array of port specs
	for _, node := range node.Content {
		err := unmarshalPortSpec(node, &spec)
		if err != nil {
			return err
		}

		p.Items = append(p.Items, spec)
	}

	return nil
}

type PortSpec interface {
}

func unmarshalPortSpec(node *yaml.Node, spec *PortSpec) error {
	switch node.Kind {

	case yaml.ScalarNode:
		var port Port
		if err := node.Decode(&port); err != nil {
			return err
		}
		*spec = &SinglePort{Port: port}
		return nil

	case yaml.MappingNode:
		var r PortRange
		if err := node.Decode(&r); err != nil {
			return err
		}
		*spec = &r
		return nil

	default:
		return fmt.Errorf("invalid port spec")
	}
}

type SinglePort struct {
	Port Port
}

type PortRange struct {
	From Port
	To   Port // Inclusive
}
