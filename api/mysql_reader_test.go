package api

import (
	"reflect"
	"testing"
)

func TestMySqlReader_LookupDomain(t *testing.T) {
	s, err := NewMySqlStorage()
	if err != nil {
		t.Fatal(err)
	}

	r, err := NewMySqlReader()
	if err != nil {
		t.Fatal(err)
	}

	z, err := s.AddZone(Zone{Name: "EU", A: []IPv4{}, AAAA: []IPv6{}})
	if err != nil {
		t.Fatal(err)
	}

	d, err := s.AddDomain(NewDomain("example.org", z))
	if err != nil {
		t.Fatal(err)
	}

	d2, err := r.LookupDomain("example.org")
	if !reflect.DeepEqual(d, d2) {
		t.Errorf("domains are not equal: %v vs %v", d2, d)
	}

	err = s.DeleteDomain(d)
	if err != nil {
		t.Fatal(err)
	}

	err = s.DeleteZone(z)
	if err != nil {
		t.Fatal(err)
	}
}
