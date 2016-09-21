package sdb

import (
	"os"
	"path/filepath"
)

// Site describes site specific configuration parameters.
type Site struct {
	SiteID int64
	// Host is a domain name with port, if port is custom.
	// Host does not include protocol.
	// Examples:
	// 		localhost:8080
	// 		noname.com
	Host string
	Path string
	// DBName is name of directory in DataRoot directory
	// and name of the database file with appended file extension:
	// Example:
	//   for datarootdir set to ~
	//   for DBName set to test
	//   name of db file derrived: ~/test/test.db
	DBName string
	Theme  string
	Title  string
	// ZonePath is a calculated property and it is not stored in DB.
	ZonePath string
}

// WebTemplateFile returns file name for 'www' directory in site specific data dir.
// 2nd return value indicates if file exists.
func (s Site) WebTemplateFile(dataDir string, fileName string) (string, bool) {
	fn := filepath.Join(dataDir, s.DBName, "t-html-s", fileName)
	if _, err := os.Stat(fn); err == nil {
		// file exists
		return fn, true
	}
	return "", false
}
