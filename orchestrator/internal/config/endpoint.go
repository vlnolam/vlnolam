package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
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
				return errors.Wrap("cidr_parse_error", err)
			}
			e.ConnectionEndpointSpec = &CIDRConnectionEndpoint{*cidr}
			return nil

		} else if ip := net.ParseIP(v); ip != nil {
			e.ConnectionEndpointSpec = &IPConnectionEndpoint{ip}
			return nil

		} else {
			return errors.New("unknown_endpoint_format", errors.Data{"value": v})
		}

	case map[string]interface{}:
		// TODO
		return errors.New("object_endpoints_not_supported", nil)

	default:
		return errors.New("invalid_endpoint_type", errors.Data{"type": fmt.Sprintf("%T", v)})
	}
}

type ConnectionEndpointSpec interface {
}

// Covers an IPv4 range
type CIDRConnectionEndpoint struct {
	net.IPNet
}

type IPConnectionEndpoint struct {
	net.IP
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
