package config

import (
	"fmt"
	"net"
	"strings"

	"go.yaml.in/yaml/v4"
)

type ConnectionEndpoint struct {
	ConnectionEndpointSpec
}

func (c ConnectionEndpoint) String() string {
	return fmt.Sprintf("%+v", c.ConnectionEndpointSpec)
}

func (e *ConnectionEndpoint) UnmarshalYAML(node *yaml.Node) error {
	var raw interface{}
	if err := node.Decode(&raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		if v == "any" {
			e.ConnectionEndpointSpec = &AnyConnectionEndpoint{}
			return nil

		} else if v == "internet" {
			e.ConnectionEndpointSpec = &InternetConnectionEndpoint{}
			return nil

		} else if strings.Contains(v, "/") {
			_, cidr, err := net.ParseCIDR(v)
			if err != nil {
				return fmt.Errorf("line %d: %w", node.Line, err)
			}
			e.ConnectionEndpointSpec = &CIDRConnectionEndpoint{IPnet: *cidr}
			return nil

		} else if ip := net.ParseIP(v); ip != nil {
			e.ConnectionEndpointSpec = &IPConnectionEndpoint{IP: ip}
			return nil

		} else {
			return fmt.Errorf("line %d: unknown endpoint format: %q", node.Line, v)
		}

	case map[string]interface{}:
		// TODO
		return fmt.Errorf("line %d: object endpoints not yet supported", node.Line)

	default:
		return fmt.Errorf("line %d: invalid endpoint type: %T", node.Line, v)
	}
}

type ConnectionEndpointSpec interface {
}

// Covers an IPv4 range
type CIDRConnectionEndpoint struct {
	IPnet net.IPNet
}

type IPConnectionEndpoint struct {
	IP net.IP
}

// Covers all IPs outside of the private ranges
type InternetConnectionEndpoint struct {
}

func (e *InternetConnectionEndpoint) String() string {
	return "internet"
}

// Covers all IPs everywhere
type AnyConnectionEndpoint struct {
}

func (e *AnyConnectionEndpoint) String() string {
	return "any"
}
