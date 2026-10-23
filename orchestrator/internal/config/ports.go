package config

import (
	"fmt"
	"strings"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
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

func (p *Ports) Validate() *errors.Error {
	for i, item := range p.Items {
		if err := item.Validate(); err != nil {
			err.PrependPath("ports/%d", i)
			return err
		}
	}
	return nil
}

type PortSpec interface {
	Validate() *errors.Error
}

func unmarshalPortSpec(node *yaml.Node, spec *PortSpec) error {
	switch node.Kind {

	case yaml.ScalarNode:
		var port Port
		if err := decodeStrict(node, &port); err != nil {
			return err
		}
		*spec = &SinglePort{Port: port}
		return nil

	case yaml.MappingNode:
		var r PortRange
		if err := decodeStrict(node, &r); err != nil {
			return err
		}
		*spec = &r
		return nil

	default:
		return errors.New("invalid_port_spec", nil)
	}
}

type SinglePort struct {
	Port Port
}

func (s *SinglePort) Validate() *errors.Error {
	if s.Port == 0 {
		return errors.New("port_zero_not_allowed", nil)
	}
	return nil
}

type PortRange struct {
	From Port
	To   Port // Inclusive
}

func (r *PortRange) Validate() *errors.Error {
	if r.From == 0 || r.To == 0 {
		return errors.New("port_zero_not_allowed", nil)
	}
	if r.From > r.To {
		return errors.New("port_range_invalid", nil)
	}
	return nil
}
