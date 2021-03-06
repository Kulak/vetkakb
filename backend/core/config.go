package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kulak/gcfg"
)

// MainSection of the configuration file.
type MainSection struct {
	// WorkingDir is used by daemon mode.
	WorkingDir string
	// SqlRoot is a parent directory of SQL scripts.
	SQLRoot string
	// WebRoot must end with slash.
	WebRoot string
	// DataRoot is a database file parent.  It not accessible by web service or not in web server subtree.
	DataRoot string
	// TemplateRoot contains HTML template and other template fles under Theme.
	TemplateRoot string
	WebEndpoint  string
	PidFileName  string
	LogFileName  string
	User         string
	// SiteURL is used to construct OAuth callback URL.
	SiteURL string
	// Client path is used in multihosted setup.
	// Path comparison is done only request URL starts with client path
	ClientPath string
}

// Configuration represents content of the configuration file.
type Configuration struct {
	Main MainSection
}

// NewConfiguration creates new configuration.
func NewConfiguration() *Configuration {
	return &Configuration{
		Main: MainSection{
			WorkingDir:   "./",
			SQLRoot:      "sql",
			WebRoot:      "www",
			DataRoot:     "data",
			TemplateRoot: "template",
			WebEndpoint:  "localhost:8080",
			PidFileName:  "/var/run/vetkakb.pid",
			LogFileName:  "/var/log/vetkakb.log",
			User:         "",
			SiteURL:      "http://localhost:8080",
			ClientPath:   "/z",
		},
	}
}

// LoadConfig loads configuration from file.
func LoadConfig(customFileName string) (*Configuration, error) {
	cfg := NewConfiguration()
	err := gcfg.ReadFileInto(cfg, customFileName)
	return cfg, err
}

// SiteDBFileName returns name of the site database.
func (c Configuration) SiteDBFileName() string {
	return c.DataFile("site.db")
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

// WebFile returns name of file in global www directory.
func (c Configuration) WebFile(fileName string) string {
	return filepath.Join(c.Main.WebRoot, fileName)
}

// TemplateThemeFile returns file name based on passed Theme name and only a file name.
func (c Configuration) TemplateThemeFile(theme, fileName string) string {
	return filepath.Join(c.Main.TemplateRoot, "theme", theme, fileName)
}

// SQLDir returns name of the directory that contains
// SQL scripts.
func (c Configuration) SQLDir(dbDir string) string {
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

	return nil
}
