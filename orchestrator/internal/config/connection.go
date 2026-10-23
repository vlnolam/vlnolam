package config

import (
	"fmt"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
	"go.yaml.in/yaml/v4"
)

type ConnectionDirection int

const (
	IngressConnection ConnectionDirection = iota
	EgressConnection
)

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

type Connection struct {
	ConnectionSpec
}

func (c Connection) String() string {
	return fmt.Sprintf("%+v", c.ConnectionSpec)
}

type header struct {
	Protocol Protocol
}

func parseConnection[body ConnectionSpec](c *Connection, node *yaml.Node) error {
	type headerAndBody struct {
		Header header `yaml:",inline"`
		Body   body   `yaml:",inline"`
	}
	data := headerAndBody{}
	err := decodeStrict(node, &data)
	c.ConnectionSpec = data.Body
	return err
}

func (c *Connection) UnmarshalYAML(node *yaml.Node) error {
	var header header
	if err := node.Decode(&header); err != nil {
		return err
	}

	switch header.Protocol {

	case "", ProtocolTCP:
		return parseConnection[*TCPConnection](c, node)

	case ProtocolUDP:
		return parseConnection[*UDPConnection](c, node)

	default:
		return fmt.Errorf("line %d: unknown protocol %q", node.Line, header.Protocol)

	}
}

type ConnectionSpec interface {
	Validate(dir ConnectionDirection) *errors.Error
}

type TCPUDPConnection struct {
	Target ConnectionEndpoint
	Source ConnectionEndpoint
	Ports  Ports
}

func (c *TCPUDPConnection) Validate(dir ConnectionDirection) *errors.Error {
	switch dir {
	case EgressConnection:
		if c.Target == (ConnectionEndpoint{}) {
			return errors.New("target_required_for_egress", nil)
		}

	case IngressConnection:
		if c.Source == (ConnectionEndpoint{}) {
			return errors.New("source_required_for_ingress", nil)
		}
	}

	if len(c.Ports.Items) == 0 {
		return errors.New("no_ports_specified", nil)
	}

	if err := c.Ports.Validate(); err != nil {
		return err
	}

	return nil
}

type TCPConnection struct {
	TCPUDPConnection `yaml:",inline"`
}

type UDPConnection struct {
	TCPUDPConnection `yaml:",inline"`
}
