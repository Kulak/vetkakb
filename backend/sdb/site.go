package sdb

// Site describes site specific configuration parameters.
type Site struct {
	SiteID int64
	// Host is a domain name with port, if port is custom.
	// Host does not include protocol.
	// Examples:
	// 		localhost:8080
	// 		noname.com
	Host   string
	Path   string
	DBName string
	Theme  string
	Title  string
	// ZonePath is a calculated property and it is not stored in DB.
	ZonePath string
}
