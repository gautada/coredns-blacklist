package blacklist

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	// "github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// parseIP calls discards any v6 zone info, before calling net.ParseIP.
func parseIP(addr string) net.IP {
        if i := strings.Index(addr, "%"); i >= 0 {
                // discard ipv6 zone
                addr = addr[0:i]
        }

        return net.ParseIP(addr)
}

// Hosts is the plugin handler
type Blacklist struct {
	Next plugin.Handler
	*ListFile

	Fall fall.F
}

// ServeDNS implements the plugin.Handle interface.
func (b Blacklist) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	answers := []dns.RR{}

	// zone := plugin.Zones(b.Origins).Matches(qname)
	// if zone == "" {
	// 	// PTR zones don't need to be specified in Origins.
	// 	if state.QType() != dns.TypePTR {
	// 		// if this doesn't match we need to fall through regardless of h.Fallthrough
	// 		return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
	// 	}
	// }

	if b.lookupDomain(qname) {
		switch state.QType() {
		case dns.TypeA:
			answers = a(qname, b.options.ttl, parseIP(b.options.ipv4))
		case dns.TypeAAAA:
			answers = aaaa(qname, b.options.ttl, parseIP(b.options.ipv6))
		}
	} else {
		if b.Fall.Through(qname) {
			return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
		} else {
			// We want to send an NXDOMAIN, but because of /etc/hosts' setup we don't have a SOA, so we make it SERVFAIL
			// to at least give an answer back to signals we're having problems resolving this.
			return dns.RcodeServerFailure, nil
		}
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the plugin.Handle interface.
func (b Blacklist) Name() string { return "blacklist" }

// a takes a slice of net.IPs and returns a slice of A RRs.
func a(zone string, ttl uint32, ip net.IP) []dns.RR {
	answers := make([]dns.RR, 1)
        r := new(dns.A)
        r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
        r.A = ip
        answers[0] = r
        return answers
}

// aaaa takes a slice of net.IPs and returns a slice of AAAA RRs.
func aaaa(zone string, ttl uint32, ip net.IP) []dns.RR {
	answers := make([]dns.RR, 1)
        r := new(dns.AAAA)
        r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl}
        r.AAAA = ip
        answers[0] = r
        return answers
}

