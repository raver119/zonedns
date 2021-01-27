package api

import (
	"reflect"
	"testing"
)

func TestMySqlStorage_Test_Zone_CRUD(t *testing.T) {
	storage, err := NewMySqlStorage()
	if err != nil {
		t.Fatal(err)
	}

	zoneA := Zone{Name: "EU", A: []IPv4{"192.168.1.12", "192.168.19.137"}, AAAA: []IPv6{"fe80::4ca0:8fff:fe70:52d"}}
	zoneB := Zone{Name: "USA", A: []IPv4{"192.168.5.17", "192.168.6.14"}, AAAA: []IPv6{"f480::5ca0:2f0f:fa91:6d0"}}

	z1, err := storage.AddZone(zoneA)
	if err != nil {
		t.Error(err)
	}

	z2, err := storage.AddZone(zoneB)
	if err != nil {
		t.Error(err)
	}

	zoneA.id = z1.id
	if !reflect.DeepEqual(zoneA, z1) {
		t.Errorf("expected: %v;\nreceived: %v", zoneA, z1)
	}

	zoneA.A = []IPv4{"172.0.0.1"}
	zoneA.AAAA = []IPv6{"fe80::4a:d4ff:fede:5f6b"}

	z1, err = storage.UpdateZone(zoneA)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(zoneA, z1) {
		t.Errorf("expected: %v;\nreceived: %v", zoneA, z1)
	}

	zones, err := storage.FetchZones()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual([]Zone{z1, z2}, zones) {
		t.Fatalf("Zones meant to be equal;\n[%v]\n[%v]", zones, []Zone{z1, z2})
	}

	err = storage.DeleteZone(z1)
	if err != nil {
		t.Error(err)
	}

	err = storage.DeleteZone(z2)
	if err != nil {
		t.Error(err)
	}
}

func TestMySqlStorage_Test_Domain_CRUD(t *testing.T) {
	storage, err := NewMySqlStorage()
	if err != nil {
		t.Fatal(err)
	}

	zone, err := storage.AddZone(Zone{Name: "EU", A: []IPv4{}, AAAA: []IPv6{}})
	if err != nil {
		t.Fatal(err)
	}

	domain := NewDomain("example.com", zone)

	d, err := storage.AddDomain(domain)
	if err != nil {
		t.Fatal(err)
	}

	domain.id = d.id
	if !reflect.DeepEqual(domain, d) {
		t.Errorf("expected: %v;\nreceived: %v", domain, d)
	}

	domain.Name = "example.org"
	domain.Txt = "ANOTHER TEXT"
	domain.ZoneID = 19

	d, err = storage.UpdateDomain(domain)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(domain, d) {
		t.Errorf("expected: %v;\nreceived: %v", domain, d)
	}

	d2, err := storage.LookupDomain(domain.Name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(d, d2) {
		t.Errorf("expected: %v;\nreceived: %v", d, d2)
	}

	err = storage.DeleteDomain(d)
	if err != nil {
		t.Error(err)
	}

	err = storage.DeleteZone(zone)
	if err != nil {
		t.Error(err)
	}
}
