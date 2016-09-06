package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kulak/gcfg"
	"github.com/kulak/sqlitemaint"
)

// MainSection of the configuration file.
type MainSection struct {
	// SqlRoot is a parent directory of SQL scripts.
	SQLRoot string
	// WebRoot must end with slash.
	WebRoot string
	// DataRoot is a database file parent.  It not accessible by web service or not in web server subtree.
	DataRoot    string
	WebEndpoint string
	PidFileName string
	LogFileName string
	User        string
}

// Configuration represents content of the configuration file.
type Configuration struct {
	Main MainSection
}

// NewConfiguration creates new configuration.
func NewConfiguration() *Configuration {
	return &Configuration{
		Main: MainSection{
			SQLRoot:     "sql",
			WebRoot:     "www",
			DataRoot:    "data",
			WebEndpoint: "localhost:8080",
			PidFileName: "/var/run/vetkakb.pid",
			LogFileName: "/var/log/vetkakb.log",
			User:        "",
		},
	}
}

// LoadConfig loads configuration from file.
func LoadConfig(cutomFileName string) (*Configuration, error) {
	fileName := "vetkakb.gcfg"
	if cutomFileName != "" {
		fileName = cutomFileName
	}
	log.Printf("Using configuration file %v", fileName)
	cfg := NewConfiguration()
	err := gcfg.ReadFileInto(cfg, fileName)
	return cfg, err
}

// EntryDBFileName returns name of the entry DB file.
func (c Configuration) EntryDBFileName() string {
	return c.DataFile("entry.db")
}

// DataFile retruns name of the file in Dataroot directory.
func (c Configuration) DataFile(fileName string) string {
	return filepath.Join(c.Main.DataRoot, fileName)
}

// WebDir returns http directory relative to WebRoot.
func (c Configuration) WebDir(dir string) http.Dir {
	// don't use filepath.Join, because it strips ending slash '/'
	d := http.Dir(c.Main.WebRoot + dir)
	log.Printf("WebDir for %s is: %s\n", dir, d)
	return d
}

// SQLDir returns name of the directory that contains
// SQL scripts.
func (c Configuration) sqlDir(dbDir string) string {
	return filepath.Join(c.Main.SQLRoot, dbDir)
}

// InitializeFilesystem creates directory tree, creates
// and updates database files.
func (c Configuration) InitializeFilesystem() (err error) {
	// create data root directory if it does not exist
	err = os.MkdirAll(c.Main.DataRoot, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create a data directory due to error: %v", err)
	}

	_, err = sqlitemaint.UpgradeSQLite(c.EntryDBFileName(), c.sqlDir("entrydb"))
	if err != nil {
		return fmt.Errorf("Failed to upgrade DB.  Error: %v", err)
	}
	return nil
}
