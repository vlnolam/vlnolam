package config

type Port uint16

type Ports struct {
	Items []PortSpec
}

type PortSpec interface {
}

type PortList struct {
	Ports []Port
}

type PortRange struct {
	From Port
	To   Port // Inclusive
}
