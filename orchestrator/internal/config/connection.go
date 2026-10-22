package config

type Connection struct {
	Spec ConnectionSpec
}

type ConnectionSpec interface {
}

type TCPConnection struct {
	Target      ConnectionEndpoint
	Source      ConnectionEndpoint
	Ports       Ports
	SourcePorts Ports
}

type UDPConnection struct {
	Target      ConnectionEndpoint
	Source      ConnectionEndpoint
	Ports       Ports
	SourcePorts Ports
}
