package zonedns

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/raver119/zonedns-go-api"
	"net"
	"regexp"
	"time"
)

var re = regexp.MustCompile("\\.$")

type ZonedNS struct {
	zones Zones
	back  api.MySqlReader
	ttl   uint32
}

func BuildZonedNS(zonesUpdateTimeoutSec int64, ttlSec uint32) (ZonedNS, error) {
	db, err := api.NewMySqlReader()
	if err != nil {
		return ZonedNS{}, err
	}

	zones, err := db.FetchZones()
	if err != nil {
		return ZonedNS{}, err
	}

	zonedns := ZonedNS{back: db, zones: BuildZones(zones), ttl: ttlSec}
	// Zones should be updated periodically
	go func() {
		for true {
			// TODO: add error handling here
			zz, err := db.FetchZones()
			if err != nil {
				panic(err)
			}

			zonedns.zones.Update(zz)

			time.Sleep(time.Duration(zonesUpdateTimeoutSec*1000) * time.Millisecond)
		}
	}()

	return zonedns, nil
}

func (z ZonedNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	var answers []dns.RR
	var zone api.Zone
	var ok bool

	dom, err := z.back.LookupDomain(re.ReplaceAllString(qname, ""))
	if err != nil {
		// NXDOMAIN
		return dns.RcodeNameError, err
	}

	if zone, ok = z.zones.Get(dom.ZoneID); !ok {
		// FIXME: try to refresh zones first?
		// NXDOMAIN as well
		return dns.RcodeNameError, err
	}

	if answers, ok = z.resolve(qname, dom, zone, state.QType()); !ok {
		return dns.RcodeServerFailure, fmt.Errorf("queryType %v failed on [%v]", state.QType(), qname)
	}

	// send result back
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	err = w.WriteMsg(m)
	return dns.RcodeSuccess, err
}

func (z ZonedNS) resolve(qname string, dom api.Domain, zone api.Zone, qType uint16) (answers []dns.RR, ok bool) {
	switch qType {
	// TODO: add reversePTR here
	case dns.TypeMX:
		answers, ok = z.resolveMX(qname, dom, zone)
	case dns.TypeTXT:
		answers, ok = z.resolveTXT(qname, dom, zone)
	case dns.TypeA:
		answers, ok = z.resolveA(qname, dom, zone)
	case dns.TypeAAAA:
		answers, ok = z.resolveA(qname, dom, zone)
	default:
		return answers, false
	}

	return
}

func (z ZonedNS) resolveMX(qname string, dom api.Domain, zone api.Zone) (answers []dns.RR, ok bool) {
	// TODO: decide something here. maybe add MX entity which will work pretty much like a zone?
	return answers, false
}

func (z ZonedNS) resolveTXT(qname string, dom api.Domain, zone api.Zone) (answers []dns.RR, ok bool) {
	if len(dom.Txt) == 0 {
		return answers, false
	}

	r := new(dns.TXT)
	r.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: z.ttl}
	r.Txt = []string{dom.Txt}

	return append(answers, r), true
}

func (z ZonedNS) resolveA(qname string, dom api.Domain, zone api.Zone) (answers []dns.RR, ok bool) {
	answers = make([]dns.RR, len(zone.A))
	for i, v := range zone.A {
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: z.ttl}
		r.A = net.ParseIP(string(v))
		answers[i] = r
	}
	ok = len(zone.A) > 0
	return
}

func (z ZonedNS) resolveAAAA(qname string, dom api.Domain, zone api.Zone) (answers []dns.RR, ok bool) {
	answers = make([]dns.RR, len(zone.AAAA))
	for i, v := range zone.AAAA {
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: z.ttl}
		r.AAAA = net.ParseIP(string(v))
		answers[i] = r
	}
	ok = len(zone.A) > 0
	return
}

func (z ZonedNS) Name() string {
	return "zonedns"
}
