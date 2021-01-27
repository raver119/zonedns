package zonedns

import (
	"github.com/miekg/dns"
	"github.com/raver119/zonedns-go-api"
	"net"
	"reflect"
	"testing"
	"time"
)

func buildZonedNS() ZonedNS {
	zz := []api.Zone{api.NewZone(1, "EU", []api.IPv4{"192.168.1.7", "192.168.2.7"}, []api.IPv6{"fe80::4a:d4ff:fede:5f6b", "fe80::f017:65ff:fe62:8e5e"})}
	return ZonedNS{zones: BuildZones(zz), ttl: 123}
}

var z = buildZonedNS()

func TestZonedNS_resolveTXT(t *testing.T) {
	zone, _ := z.zones.Get(1)
	domain := api.NewDomain("example.org", zone)
	domain.Txt = "TEXT RECORD"

	answers, ok := z.resolveTXT("example.org", domain, zone)
	if !ok {
		t.Fatalf("request should've been resolved properly")
	}

	if len(answers) != 1 {
		t.Fatalf("there should be exactly 1 response")
	}

	if answers[0].Header().Rrtype != dns.TypeTXT {
		t.Errorf("TypeTXT expected, but got %v instead", answers[0].Header().Rrtype)
	}

	if answers[0].Header().Ttl != 123 {
		t.Errorf("TTL of 123 expected, but got %v instead", answers[0].Header().Ttl)
	}
}

func TestZonedNS_resolveA(t *testing.T) {
	zone, _ := z.zones.Get(1)
	domain := api.NewDomain("example.org", zone)

	answers, ok := z.resolveA("example.org", domain, zone)
	if !ok {
		t.Fatalf("request should've been resolved properly")
	}

	if len(answers) != 2 {
		t.Fatalf("there should be exactly 2 responses")
	}

	if answers[0].Header().Rrtype != dns.TypeA {
		t.Errorf("TypeAAAA expected, but got %v instead", answers[0].Header().Rrtype)
	}

	if answers[0].Header().Ttl != 123 {
		t.Errorf("TTL of 123 expected, but got %v instead", answers[0].Header().Ttl)
	}

	if !reflect.DeepEqual(answers[0].(*dns.A).A, net.ParseIP(string(zone.A[0]))) {
		t.Fatalf("IPs should match")
	}

	if !reflect.DeepEqual(answers[1].(*dns.A).A, net.ParseIP(string(zone.A[1]))) {
		t.Fatalf("IPs should match")
	}
}

func TestZonedNS_resolveAAAA(t *testing.T) {
	zone, _ := z.zones.Get(1)
	domain := api.NewDomain("example.org", zone)

	answers, ok := z.resolveAAAA("example.org", domain, zone)
	if !ok {
		t.Fatalf("request should've been resolved properly")
	}

	if len(answers) != 2 {
		t.Fatalf("there should be exactly 2 responses")
	}

	if answers[0].Header().Rrtype != dns.TypeAAAA {
		t.Errorf("TypeAAAA expected, but got %v instead", answers[0].Header().Rrtype)
	}

	if answers[0].Header().Ttl != 123 {
		t.Errorf("TTL of 123 expected, but got %v instead", answers[0].Header().Ttl)
	}

	if !reflect.DeepEqual(answers[0].(*dns.AAAA).AAAA, net.ParseIP(string(zone.AAAA[0]))) {
		t.Fatalf("IPs should match")
	}

	if !reflect.DeepEqual(answers[1].(*dns.AAAA).AAAA, net.ParseIP(string(zone.AAAA[1]))) {
		t.Fatalf("IPs should match")
	}
}

func TestBuildZonedNS(t *testing.T) {
	storage, err := api.NewMySqlStorage()
	if err != nil {
		t.Fatal(err)
	}

	zone, err := storage.AddZone(api.Zone{Name: "EU", A: []api.IPv4{"192.168.1.7"}})
	if err != nil {
		t.Fatal(err)
	}

	// short update cycle
	ns, err := BuildZonedNS(1, 600)

	if err != nil {
		t.Fatal(err)
	}

	if !ns.zones.Has(zone.Id()) {
		t.Errorf("can't see zone %v in zones", zone.Id())
	}

	zone2, err := storage.AddZone(api.Zone{Name: "US", A: []api.IPv4{"192.168.2.2"}})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	if !ns.zones.Has(zone.Id()) {
		t.Errorf("can't see zone %v in zones", zone.Id())
	}
	if !ns.zones.Has(zone2.Id()) {
		t.Errorf("can't see new zone %v in zones", zone.Id())
	}

	err = storage.DeleteZone(zone)
	if err != nil {
		t.Error(err)
	}
	err = storage.DeleteZone(zone2)
	if err != nil {
		t.Error(err)
	}
}
