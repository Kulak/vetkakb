package sdb

import (
	"database/sql"
	"fmt"

	"github.com/kulak/sqlitemaint"
)

// SiteDB represents content of site database.
type SiteDB struct {
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

// Upgrade creates and upgrades Site datbase.
func (sdb *SiteDB) Upgrade() (err error) {
	_, err = sqlitemaint.UpgradeSQLite(sdb.dbFileName, sdb.sqlDir)
	if err != nil {
		return fmt.Errorf("Failed to upgrade site DB %s.  Error: %v", sdb.dbFileName, err)
	}
	return
}

// Open opens DB connection.
func (sdb *SiteDB) Open() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", sdb.dbFileName)
	if err != nil {
		return db, fmt.Errorf("Failed to open database %v. Error: %v", sdb.dbFileName, err)
	}
	return
}

// GetSite loads site by hostname and path.
func (sdb *SiteDB) GetSite(host, path string) (site *Site, err error) {
	var db *sql.DB
	db, err = sdb.Open()
	if err != nil {
		return
	}
	defer db.Close()

	query := `
	SELECT s.siteID, s.dbname, s.theme, s.title from site s
	where s.host = ? and s.path = ?;
	`
	site = &Site{Host: host, Path: path}
	err = db.QueryRow(query, host, path).Scan(
		&site.SiteID, &site.DBName, &site.Theme, &site.Title)
	// log.Printf("Loaded site id %v with dbName %s for host %s, path %s.  Err: %v",
	// 	site.SiteID, site.DBName, host, path, err)
	return
}

// GetSiteByID loads site by hostname and path.
func (sdb *SiteDB) GetSiteByID(siteIDStr string) (s *Site, err error) {
	var db *sql.DB
	db, err = sdb.Open()
	if err != nil {
		return
	}
	defer db.Close()

	query := `
	SELECT siteID, host, path, dbname, theme, title from site
	where siteID = ?;
	`
	s = &Site{}
	err = db.QueryRow(query, siteIDStr).Scan(
		&s.SiteID, &s.Host, &s.Path, &s.DBName, &s.Theme, &s.Title)
	// log.Printf("Loaded site id %v with dbName %s for host %s, path %s.  Err: %v",
	// 	s.SiteID, s.DBName, s.Host, s.Path, err)
	return
}

// All loads all sites.
func (sdb *SiteDB) All() (sites []*Site, err error) {
	var db *sql.DB
	db, err = sdb.Open()
	if err != nil {
		return
	}
	defer db.Close()

	query := `SELECT siteID, host, path, dbname, theme, title from site`
	var rows *sql.Rows
	rows, err = db.Query(query)
	if err != nil {
		return
	}
	for rows.Next() {
		// NOOTE: keep variable in the loop.
		// If outside, then pointer in append points to the same memory location
		// This way memory location is new for each loop iteration!
		s := Site{}
		err = rows.Scan(&s.SiteID, &s.Host, &s.Path, &s.DBName, &s.Theme, &s.Title)
		if err != nil {
			return
		}
		sites = append(sites, &s)
	}
	return
}
