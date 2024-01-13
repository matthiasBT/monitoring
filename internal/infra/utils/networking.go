package utils

import "net"

func ParseSubnet(rawSubnet string) *net.IPNet {
	if rawSubnet == "" {
		panic("Empty subnet string")
	}
	_, subnet, err := net.ParseCIDR(rawSubnet)
	if err != nil {
		panic("Invalid CIDR notation")
	}
	return subnet
}
