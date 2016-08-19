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

// AppRoot is a directory that's a parent of www directory.
type MainSection struct {
	AppRoot     string
	DataRoot    string
	WebEndpoint string
	PidFileName string
	LogFileName string
	User        string
}

type Configuration struct {
	Main MainSection
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Main: MainSection{
			AppRoot:     "",
			DataRoot:    "data",
			WebEndpoint: "localhost:8080",
			PidFileName: "/var/run/vetkakb.pid",
			LogFileName: "/var/log/vetkakb.log",
			User:        "",
		},
	}
}

func LoadConfig(cutomFileName string) (*Configuration, error) {
	file_name := "vetkakb.gcfg"
	if cutomFileName != "" {
		file_name = cutomFileName
	}
	log.Printf("Using configuration file %v", file_name)
	cfg := NewConfiguration()
	err := gcfg.ReadFileInto(cfg, file_name)
	return cfg, err
}

// Returns http directory relative to AppRoot.
func (self Configuration) WebDir(dir string) http.Dir {
	// don't use filepath.Join, because it strips ending slash '/'
	d := http.Dir(self.Main.AppRoot + dir)
	log.Printf("WebDir for %s is: %s\n", dir, d)
	return d
}

func (self Configuration) VetkaDBFileName() string {
	return filepath.Join(self.Main.DataRoot, "vetka.db")
}

func (self Configuration) SqlDir(dbDir string) string {
	return filepath.Join(self.Main.AppRoot, "sql", dbDir)
}

func (c Configuration) InitializeFilesystem() (err error) {
	// create data root directory if it does not exist
	err = os.MkdirAll(c.Main.DataRoot, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create a data directory due to error: %v", err)
	}

	_, err = sqlitemaint.UpgradeSQLite(c.VetkaDBFileName(), c.SqlDir("vetka"))
	if err != nil {
		return fmt.Errorf("Failed to upgrade DB.  Error: %v", err)
	}
	return nil
}
