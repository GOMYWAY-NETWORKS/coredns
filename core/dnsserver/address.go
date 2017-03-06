package dnsserver

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type zoneAddr struct {
	Zone     string
	Port     string
	Protocol string // dns, tls or grpc (not implemented)
}

// String return z.Zone + ":" + z.Port as a string.
func (z zoneAddr) String() string { return z.Protocol + "://" + z.Zone + ":" + z.Port }

// Protocol returns the protocol of the string s
func Protocol(s string) string {
	switch {
	case strings.HasPrefix(s, ProtoTLS+"://"):
		return ProtoTLS
	case strings.HasPrefix(s, ProtoDNS+"://"):
		return ProtoDNS
	}
	return ProtoDNS
}

// normalizeZone parses an zone string into a structured format with separate
// host, and port portions, as well as the original input string.
//
// TODO(miek): possibly move this to middleware/normalize.go
func normalizeZone(str string) (zoneAddr, error) {
	var err error

	proto := ProtoDNS

	switch {
	case strings.HasPrefix(str, ProtoTLS+"://"):
		proto = ProtoTLS
		str = str[len(ProtoTLS+"://"):]
	case strings.HasPrefix(str, ProtoDNS+"://"):
		proto = ProtoDNS
		str = str[len(ProtoDNS+"://"):]
		// error if nothing matches? TODO
	}

	host, port, err := net.SplitHostPort(str)
	if err != nil {
		host, port, err = net.SplitHostPort(str + ":")
		// no error check here; return err at end of function
	}

	if len(host) > 255 {
		return zoneAddr{}, fmt.Errorf("specified zone is too long: %d > 255", len(host))
	}
	_, d := dns.IsDomainName(host)
	if !d {
		return zoneAddr{}, fmt.Errorf("zone is not a valid domain name: %s", host)
	}

	if port == "" {
		if proto == ProtoDNS {
			port = Port
		}
		if proto == ProtoTLS {
			port = TLSPort
		}
	}

	return zoneAddr{Zone: strings.ToLower(dns.Fqdn(host)), Port: port, Protocol: proto}, err
}

// Support protocols
const (
	ProtoDNS = "dns"
	ProtoTLS = "tls"

//	ProtogRPC = "grpc"
)
