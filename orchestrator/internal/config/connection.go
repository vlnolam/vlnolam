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

func (c *Connection) UnmarshalYAML(node *yaml.Node) error {
	var header ConnectionHeader
	if err := node.Decode(&header); err != nil {
		return err
	}

	switch header.Protocol {

	case "", ProtocolTCP:
		conn := &TCPConnection{}
		// Fill in for the "" case
		// Must come after decode bcs it could overwrite it (by using protocol: "")
		defer conn.setProtocol(ProtocolTCP)
		c.ConnectionSpec = conn

	case ProtocolUDP:
		c.ConnectionSpec = &UDPConnection{}

	default:
		return fmt.Errorf("line %d: unknown protocol %q", node.Line, header.Protocol)

	}

	return decodeStrict(node, c.ConnectionSpec)
}

type ConnectionSpec interface {
	Validate(dir ConnectionDirection) *errors.Error
}

type ConnectionHeader struct {
	Protocol Protocol
}

func (h *ConnectionHeader) setProtocol(p Protocol) {
	h.Protocol = p
}

type TCPUDPConnection struct {
	ConnectionHeader `yaml:",inline"`
	Target           ConnectionEndpoint
	Source           ConnectionEndpoint
	Ports            any // TODO
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

	/*if len(c.Ports.Items) == 0 {
		return errors.New(errors.ErrConfigValidation, "no ports specified", nil)
	}*/

	return nil
}

type TCPConnection struct {
	TCPUDPConnection `yaml:",inline"`
}

type UDPConnection struct {
	TCPUDPConnection `yaml:",inline"`
}
