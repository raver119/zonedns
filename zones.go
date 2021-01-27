package zonedns

import (
	api "github.com/raver119/zonedns-go-api"
	"reflect"
	"sync"
)

type Zones struct {
	// since reads/writes ratio is going to be HUGE, regular map with a single lock is ok
	m map[int64]api.Zone
	l *sync.Mutex

	z []api.Zone
}

func contains(haystack []api.Zone, needle api.Zone) bool {
	for _, v := range haystack {
		if v.Id() == needle.Id() {
			return true
		}
	}

	return false
}

func process(existing []api.Zone, updates []api.Zone) (r []api.Zone, u []api.Zone) {
	// search for removed first
	for _, v := range existing {
		if !contains(updates, v) {
			r = append(r, v)
		}
	}

	return r, updates
}

func BuildZones(zones []api.Zone) Zones {
	m := make(map[int64]api.Zone)
	z := make([]api.Zone, 0)

	for _, v := range zones {
		m[v.Id()] = v
		c := v
		z = append(z, c)
	}

	return Zones{m: m, l: new(sync.Mutex), z: z}
}

func (z *Zones) Update(zones []api.Zone) {
	z.l.Lock()
	defer z.l.Unlock()

	// shortcut for early exit
	if reflect.DeepEqual(z.z, zones) {
		return
	}

	removed, updated := process(z.z, zones)

	// remove unwanted zones first
	for _, v := range removed {
		if _, ok := z.m[v.Id()]; ok {
			delete(z.m, v.Id())
		}
	}

	// update others
	for _, v := range updated {
		z.m[v.Id()] = v
	}

	// update reference
	z.z = zones
}

func (z *Zones) Get(zoneId int64) (zone api.Zone, ok bool) {
	z.l.Lock()
	defer z.l.Unlock()

	zone, ok = z.m[zoneId]
	return
}

func (z *Zones) Has(zoneId int64) (ok bool) {
	z.l.Lock()
	defer z.l.Unlock()

	_, ok = z.m[zoneId]
	return
}
