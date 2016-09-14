package sdb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kulak/sqlitemaint"
)

// SiteDB represents content of site database.
type SiteDB struct {
	db         *sql.DB
	dbFileName string
	sqlDir     string
}

// NewSiteDB creates site DB service.
func NewSiteDB(dbFileName, sqlDir string) *SiteDB {
	return &SiteDB{
		dbFileName: dbFileName,
		sqlDir:     sqlDir,
	}
}

// Open opens DB connection.
func (sdb *SiteDB) Open() error {
	if sdb.db != nil {
		return nil
	}

	_, err := sqlitemaint.UpgradeSQLite(sdb.dbFileName, sdb.sqlDir)
	if err != nil {
		return fmt.Errorf("Failed to upgrade site DB %s.  Error: %v", sdb.dbFileName, err)
	}

	db, err := sql.Open("sqlite3", sdb.dbFileName)
	if err != nil {
		return fmt.Errorf("Failed to open database %v. Error: %v", sdb.dbFileName, err)
	}
	sdb.db = db
	return nil
}

// Close closes DB connection if it is open.
func (sdb *SiteDB) Close() {
	if sdb.db != nil {
		sdb.db.Close()
		sdb.db = nil
	}
}

// GetSite loads site by hostname and path.
func (sdb *SiteDB) GetSite(host, path string) (site *Site, err error) {
	if sdb.db == nil {
		return site, fmt.Errorf("Database connection is closed.")
	}
	query := `
	SELECT s.siteID, s.dbname, s.theme, s.title from site s
	where s.host = ? and s.path = ?;
	`
	site = &Site{Host: host, Path: path}
	err = sdb.db.QueryRow(query, host, path).Scan(
		&site.SiteID, &site.DBName, &site.Theme, &site.Title)
	log.Printf("Loaded site id %v with dbName %s for host %s, path %s.  Err: %v",
		site.SiteID, site.DBName, host, path, err)
	return
}

// GetSite loads site by hostname and path.
func (sdb *SiteDB) GetSiteByID(siteIDStr string) (s *Site, err error) {
	if sdb.db == nil {
		return s, fmt.Errorf("Database connection is closed.")
	}
	query := `
	SELECT siteID, host, path, dbname, theme, title from site
	where siteID = ?;
	`
	s = &Site{}
	err = sdb.db.QueryRow(query, siteIDStr).Scan(
		&s.SiteID, &s.Host, &s.Path, &s.DBName, &s.Theme, &s.Title)
	log.Printf("Loaded site id %v with dbName %s for host %s, path %s.  Err: %v",
		s.SiteID, s.DBName, s.Host, s.Path, err)
	return
}

// All loads all sites.
func (sdb *SiteDB) All() (sites []*Site, err error) {
	if sdb.db == nil {
		return sites, fmt.Errorf("Database connection is closed.")
	}
	query := `SELECT siteID, host, path, dbname, theme, title from site`
	var rows *sql.Rows
	rows, err = sdb.db.Query(query)
	if err != nil {
		return
	}
	s := &Site{}
	for rows.Next() {
		err = rows.Scan(&s.SiteID, &s.Host, &s.Path, &s.DBName, &s.Theme, &s.Title)
		if err != nil {
			return
		}
		sites = append(sites, s)
	}
	return
}
