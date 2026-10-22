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

type Connection struct {
	ConnectionSpec
}

func (c Connection) String() string {
	return fmt.Sprintf("Connection{Spec:%+v}", c.ConnectionSpec)
}

func (c *Connection) UnmarshalYAML(node *yaml.Node) error {
	var header ConnectionHeader
	if err := node.Decode(&header); err != nil {
		return err
	}

	switch header.Protocol {
	case "", "tcp":
		c.ConnectionSpec = &TCPConnection{}
	case "udp":
		c.ConnectionSpec = &UDPConnection{}
	default:
		return fmt.Errorf("line %d: unknown protocol %q", node.Line, header.Protocol)
	}

	return node.Decode(c.ConnectionSpec)
}

type ConnectionSpec interface {
	Validate(dir ConnectionDirection) *errors.Error
}

type ConnectionHeader struct {
	Protocol string
}

type TCPUDPConnection struct {
	ConnectionHeader
	Target ConnectionEndpoint
	Source ConnectionEndpoint
	Ports  Ports
}

func (c *TCPUDPConnection) Validate(dir ConnectionDirection) *errors.Error {
	switch dir {
	case EgressConnection:
		if c.Target == nil {
			return errors.New(errors.ErrConfigValidation, "target required for egress connections")
		}

	case IngressConnection:
		if c.Source == nil {
			return errors.New(errors.ErrConfigValidation, "source required for ingress connections")
		}
	}

	if len(c.Ports.Items) == 0 {
		return errors.New(errors.ErrConfigValidation, "no ports specified")
	}

	return nil
}

type TCPConnection struct {
	TCPUDPConnection
}

type UDPConnection struct {
	TCPUDPConnection
}
