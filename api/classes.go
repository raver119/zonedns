package api

import (
	"net"
)

type IPv4 string
type IPv6 string

type Zone struct {
	id   int64
	Name string
	A    []IPv4
	AAAA []IPv6
}

type Domain struct {
	id     int64
	Name   string
	ZoneID int64
	Txt    string
}

func validateIPv4(v IPv4) bool {
	return net.ParseIP(string(v)) != nil
}

func validateIPv6(v IPv6) bool {
	return net.ParseIP(string(v)) != nil
}

func NewDomain(domain string, zone Zone) Domain {
	return Domain{
		Name:   domain,
		ZoneID: zone.id,
		Txt:    "",
	}
}
