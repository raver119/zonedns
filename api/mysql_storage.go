package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type MySqlStorage struct {
	db *sql.DB
	r  MySqlReader
}

func GetEnvOrDefault(variable string, defaultValue string) string {
	if val, ok := os.LookupEnv(variable); ok {
		return val
	} else {
		return defaultValue
	}
}

func NewMySqlStorage() (MySqlStorage, error) {
	username := GetEnvOrDefault("DB_USERNAME", "")
	password := GetEnvOrDefault("DB_PASSWORD", "")
	hostname := GetEnvOrDefault("DB_HOST", "localhost")
	database := GetEnvOrDefault("DB_NAME", "zonedns")

	return BuildMySqlStorage(username, password, hostname, database)
}

func buildMySqlConnection(username string, password string, hostname string, database string) (db *sql.DB, err error) {
	address := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", username, password, hostname, database)
	fmt.Printf("Connecting to MySQL server: %v:%v@tcp(%v)/%v?charset=utf8\n", username, "******", hostname, database)

	db, err = sql.Open("mysql", address)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return
}

func BuildMySqlStorage(username string, password string, hostname string, database string) (s MySqlStorage, err error) {
	db, err := buildMySqlConnection(username, password, hostname, database)
	if err != nil {
		return
	}

	s = MySqlStorage{db, MySqlReader{db}}
	err = s.init()
	if err != nil {
		return MySqlStorage{}, err
	}

	return s, err
}

func (m MySqlStorage) FetchZones() (z []Zone, err error) {
	return m.r.FetchZones()
}

func encodeAA(zone Zone) (a []byte, aaaa []byte, err error) {
	a, err = json.Marshal(zone.A)
	if err != nil {
		return
	}

	aaaa, err = json.Marshal(zone.AAAA)
	if err != nil {
		return
	}

	return
}

func (m MySqlStorage) validateZone(z Zone) (err error) {
	for _, v := range z.A {
		if !validateIPv4(v) {
			return fmt.Errorf("invalid IPv4 address provided in Zone A: %v", v)
		}
	}

	for _, v := range z.AAAA {
		if !validateIPv6(v) {
			return fmt.Errorf("invalid IPv6 address provided in Zone AAAA: %v", v)
		}
	}

	return nil
}

func (m MySqlStorage) init() (err error) {
	_, err = m.db.Exec("create table IF NOT EXISTS ZNS_zones\n(\n    z_id   bigint auto_increment\n        primary key,\n    z_name varchar(128) default '' not null,\n    z_a    json                    null,\n    z_aaaa json                    null\n);\n")
	if err != nil {
		return err
	}

	_, err = m.db.Exec("create table IF NOT EXISTS ZNS_domains\n(\n    d_id   bigint auto_increment\n        primary key,\n    d_zid  bigint       not null,\n    d_name varchar(255) not null,\n    d_mx   json         not null,\n    d_txt  text         null,\n    constraint ZNS_domains_d_name_uindex\n        unique (d_name)\n);")
	if err != nil {
		return err
	}

	return
}

func (m MySqlStorage) AddZone(zone Zone) (z Zone, err error) {
	a, aaaa, err := encodeAA(zone)
	if err != nil {
		return Zone{}, err
	}

	err = m.validateZone(zone)
	if err != nil {
		return Zone{}, err
	}

	res, err := m.db.Exec("INSERT INTO ZNS_zones(z_name, z_a, z_aaaa) values (?, ?, ?)", zone.Name, a, aaaa)
	if err != nil {
		return Zone{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Zone{}, err
	}

	return m.r.fetchZoneById(id)
}

func (m MySqlStorage) UpdateZone(zone Zone) (z Zone, err error) {
	a, aaaa, err := encodeAA(zone)
	if err != nil {
		return Zone{}, err
	}

	err = m.validateZone(zone)
	if err != nil {
		return Zone{}, err
	}

	_, err = m.db.Exec("UPDATE ZNS_zones SET z_name=?, z_a=?, z_aaaa=? WHERE z_id=?", zone.Name, a, aaaa, zone.id)
	if err != nil {
		return Zone{}, err
	}

	return m.r.fetchZoneById(zone.id)
}

func (m MySqlStorage) DeleteZone(zone Zone) (err error) {
	return m.DeleteZoneById(zone.id)
}

func (m MySqlStorage) DeleteZoneById(zoneId int64) (err error) {
	_, err = m.db.Exec("DELETE from ZNS_zones WHERE z_id=? LIMIT 1", zoneId)
	return
}

func (m MySqlStorage) LookupDomain(domain string) (d Domain, err error) {
	return m.r.LookupDomain(domain)
}

/*
	This method adds specified domain
*/
func (m MySqlStorage) AddDomainAsString(domain string, zoneId int64) (d Domain, err error) {
	return m.AddDomain(Domain{Name: domain, ZoneID: zoneId})
}

/*
	This method adds specified domain to the database
*/
func (m MySqlStorage) AddDomain(domain Domain) (d Domain, err error) {
	_, err = m.r.fetchZoneById(domain.ZoneID)
	if err != nil {
		err = fmt.Errorf("unable to find specified zone %v", domain.ZoneID)
		return
	}

	res, err := m.db.Exec("INSERT INTO ZNS_domains(d_zid, d_name, d_txt, d_mx) VALUES (?, ?, ?, ?)", domain.ZoneID, domain.Name, domain.Txt, "{}")
	if err != nil {
		return Domain{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Domain{}, err
	}

	return m.r.fetchDomainById(id)
}

/*
	This method updates specified domain in the database.
*/
func (m MySqlStorage) UpdateDomain(d Domain) (dom Domain, err error) {
	_, err = m.db.Exec("UPDATE ZNS_domains SET d_zid=?, d_name=?, d_txt=? WHERE d_id=? LIMIT 1", d.ZoneID, d.Name, d.Txt, d.id)
	if err != nil {
		return Domain{}, err
	}

	return m.r.fetchDomainById(d.id)
}

/*
	This method allows to delete specified domain from the database
*/
func (m MySqlStorage) DeleteDomain(domain Domain) (err error) {
	return m.DeleteDomainById(domain.id)
}

/*
	This method allows to delete domain with specified ID from the database
*/
func (m MySqlStorage) DeleteDomainById(domainId int64) (err error) {
	_, err = m.db.Exec("DELETE FROM ZNS_domains WHERE d_id=? LIMIT 1", domainId)
	return
}
