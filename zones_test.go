package zonedns

import (
	api "github.com/raver119/zonedns-go-api"
	"reflect"
	"testing"
)

func TestZones_Has(t *testing.T) {
	zz := []api.Zone{api.NewZone4(0, "EU-1", []api.IPv4{"192.168.1.7", "192.168.1.9"}), api.NewZone4(1, "EU-2", []api.IPv4{"192.168.99.18", "192.168.99.182"})}
	zones := BuildZones(zz)

	if !zones.Has(0) {
		t.Errorf("expected to have zone 0")
	}

	if zones.Has(119) {
		t.Errorf("expected to have no zone 119")
	}
}

func TestZones_Get(t *testing.T) {
	zz := []api.Zone{api.NewZone4(0, "EU-1", []api.IPv4{"192.168.1.7", "192.168.1.9"}), api.NewZone4(1, "EU-2", []api.IPv4{"192.168.99.18", "192.168.99.182"})}
	zones := BuildZones(zz)

	if z, ok := zones.Get(1); ok {
		if !reflect.DeepEqual(zz[1], z) {
			t.Errorf("zones supposed to be equal:\n %v\n %v", zz[1], z)
		}
	} else {
		t.Errorf("expected to have zone 1")
	}
}

func TestZones_Update(t *testing.T) {
	zz1 := []api.Zone{api.NewZone4(0, "EU-1", []api.IPv4{"192.168.1.7", "192.168.1.9"}), api.NewZone4(1, "EU-2", []api.IPv4{"192.168.99.18", "192.168.99.182"})}
	zz2 := []api.Zone{api.NewZone4(1, "EU-2", []api.IPv4{"192.168.99.18", "192.168.99.182"}), api.NewZone4(2, "EU-3", []api.IPv4{"10.1.2.3"})}

	z := BuildZones(zz1)
	z.Update(zz2)

	if !z.Has(2) {
		t.Errorf("expected to have zone 2")
	}

	if z.Has(0) {
		t.Errorf("expected to have no zone 0")
	}
}
