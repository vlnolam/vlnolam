package config

type ConnectionEndpoint interface {
}

// Covers an IPv4 range
type IPv4ConnectionEndpoint struct {
	Ip   uint32
	Mask uint32
}

// Covers all IPs outside of the private ranges
type InternetConnectionEndpoint struct {
}
