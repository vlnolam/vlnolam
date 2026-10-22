package config

import (
	"fmt"

	"go.yaml.in/yaml/v4"
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
}

type ConnectionHeader struct {
	Protocol string `yaml:"protocol"`
}

type TCPUDPConnection struct {
	ConnectionHeader
	Target      ConnectionEndpoint
	Source      ConnectionEndpoint
	Ports       Ports
	SourcePorts Ports
}

type TCPConnection struct {
	TCPUDPConnection
}

type UDPConnection struct {
	TCPUDPConnection
}
