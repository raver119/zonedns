package api

type ZoneStorage interface {
	FetchZones() (z []Zone, err error)
	AddZone(zone Zone) (z Zone, err error)
	UpdateZone(zone Zone) (z Zone, err error)
	DeleteZone(zone Zone) (err error)
	DeleteZoneById(zoneId int64) (err error)

	LookupDomain(domain string) (d Domain, err error)
	AddDomainAsString(domain string, zoneId int64) (d Domain, err error)
	AddDomain(domain Domain) (d Domain, err error)
	UpdateDomain(domain Domain) (d Domain, err error)
	DeleteDomain(domain Domain) (err error)
	DeleteDomainById(domainId int64) (err error)
}

type ZoneReader interface {
	FetchZones() (z []Zone, err error)
	LookupDomain(domain string) (d Domain, err error)
}
