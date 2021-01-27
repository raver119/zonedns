package api

import (
	"database/sql"
	"encoding/json"
)
import _ "github.com/go-sql-driver/mysql"

type MySqlReader struct {
	db *sql.DB
}

func NewMySqlReader() (MySqlReader, error) {
	username := GetEnvOrDefault("DB_USERNAME", "")
	password := GetEnvOrDefault("DB_PASSWORD", "")
	hostname := GetEnvOrDefault("DB_HOST", "localhost")
	database := GetEnvOrDefault("DB_NAME", "zonedns")

	return BuildMySqlReader(username, password, hostname, database)
}

func BuildMySqlReader(username string, password string, hostname string, database string) (MySqlReader, error) {
	db, err := buildMySqlConnection(username, password, hostname, database)
	return MySqlReader{db: db}, err
}

func (m *MySqlReader) fetchZoneById(zoneId int64) (z Zone, err error) {
	row := m.db.QueryRow("SELECT z_id, z_name, z_a, z_aaaa FROM ZNS_zones WHERE z_id=?", zoneId)

	var a []byte
	var aaaa []byte

	err = row.Scan(&z.id, &z.Name, &a, &aaaa)
	if err != nil {
		return
	}

	err = json.Unmarshal(a, &z.A)
	if err != nil {
		return
	}

	err = json.Unmarshal(aaaa, &z.AAAA)
	if err != nil {
		return
	}

	return
}

func (m *MySqlReader) fetchDomainById(domainId int64) (d Domain, err error) {
	// TODO: add MX here?
	row := m.db.QueryRow("SELECT d_id, d_zid, d_name, d_txt FROM ZNS_domains WHERE d_id=?", domainId)
	err = row.Scan(&d.id, &d.ZoneID, &d.Name, &d.Txt)
	return
}

/*
	This methods retrieves list of Zones from the database server
*/
func (m *MySqlReader) FetchZones() (z []Zone, err error) {
	rows, err := m.db.Query("SELECT z_id, z_name, z_a, z_aaaa from ZNS_zones ORDER BY z_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var a []byte
	var aaaa []byte

	for rows.Next() {
		var zone Zone

		err = rows.Scan(&zone.id, &zone.Name, &a, &aaaa)
		if err != nil {
			return
		}

		err = json.Unmarshal(a, &zone.A)
		if err != nil {
			return
		}

		err = json.Unmarshal(aaaa, &zone.AAAA)
		if err != nil {
			return
		}

		z = append(z, zone)
	}

	return
}

/*
	This method retrieves Zone id from the database server. This ZoneID can be converted to A or AAAA record.
*/
func (m *MySqlReader) LookupDomain(domain string) (d Domain, err error) {
	row := m.db.QueryRow("SELECT d_id, d_zid, d_name, d_txt from ZNS_domains WHERE d_name=? LIMIT 1", domain)
	err = row.Scan(&d.id, &d.ZoneID, &d.Name, &d.Txt)
	return
}
