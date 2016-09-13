package sdb

// Site describes site specific configuration parameters.
type Site struct {
	SiteID int64
	Host   string
	Path   string
	DBName string
	Theme  string
}
