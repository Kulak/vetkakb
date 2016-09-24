package main

// Building Code
//
// If built with
// 		go build
// then output name is taken from parent directory name (vetkakb).
//
// If built with
//		go build vetkakb.go
// then output name is taken from file name (vetkakb).
//
// That's why name of the main file and directory match.

import (
	"flag"
	"log"
	"log/syslog"
	"net/http"
	"os/user"
	"strconv"
	"syscall"

	"horse.lan.gnezdovi.com/vetkakb/backend/core"
	"horse.lan.gnezdovi.com/vetkakb/backend/edb"
	"horse.lan.gnezdovi.com/vetkakb/backend/sdb"
	"horse.lan.gnezdovi.com/vetkakb/backend/vetka"

	"github.com/sevlyar/go-daemon"
)

// Example:
//    vetkakb -d false
//
func main() {
	var configFile string
	var consoleMode bool
	flag.StringVar(&configFile, "cf", "/usr/local/etc/vetkakb.ini", "Configuration file name")
	flag.BoolVar(&consoleMode, "c", false, "Console mode (non-daemon)")
	flag.Parse()
	log.Println("*** Starting rashodi ***")
	log.Printf("* Config file:  %s", configFile)
	log.Printf("* Console mode: %v", consoleMode)

	conf, err := core.LoadConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if consoleMode {
		// blocks
		childRun(conf)
	} else {
		// blocks
		deamonRun(conf)
	}
}

func deamonRun(conf *core.Configuration) {
	var err error

	logwriter, err := syslog.New(syslog.LOG_NOTICE, "vetkakb")
	if err != nil {
		log.Fatalf("Failed to start syslog: %s", err)
	}
	log.SetOutput(logwriter)
	log.SetFlags(0)

	context := &daemon.Context{
		PidFileName: conf.Main.PidFileName,
		PidFilePerm: 0644,
		LogFileName: conf.Main.LogFileName,
		LogFilePerm: 0640,
		Credential:  &syscall.Credential{},
		WorkDir:     conf.Main.WorkingDir,
		//Umask:       027,
	}
	if conf.Main.User != "" {
		// lookup user
		log.Printf("Looking up user %s", conf.Main.User)
		var u *user.User
		u, err = user.Lookup(conf.Main.User)
		if err != nil {
			log.Fatalf("Failed to lookup user.  Error: %v", err)
		}
		var uid int
		uid, err = strconv.Atoi(u.Uid)
		if err != nil {
			log.Fatalf("Failed to convert userid to a number.  Error: %v", err)
		}
		context.Credential.Uid = uint32(uid)
	}

	child, err := context.Reborn()
	if err != nil {
		log.Fatalf("Failed to fork rashodi.  Error: %v", err)
	}
	if child != nil {
		// parent does nothing with its child
		log.Println("Master process ended.")
		return
	}

	// child code beyond this point
	log.Println("Child process continues.")
	defer context.Release()
	childRun(conf)
}

func childRun(conf *core.Configuration) {
	var err error
	err = conf.InitializeFilesystem()
	if err != nil {
		log.Fatalf("Failed to InitializeFilesystem.  Error: %v", err)
	}

	ts := edb.NewTypeService()
	ts.Initialize()

	// edb := edb.NewEntryDB(conf.SQLDir("entrydb"), c.Main.DataRoot, ,  ts)
	// edb.Open()

	sdb := sdb.NewSiteDB(conf.SiteDBFileName(), conf.SQLDir("sitedb"))
	sdb.Open()

	log.Println("Startign web service")
	ws := vetka.NewWebSvc(conf, sdb, ts)

	// initialized listed sites
	sites, err := sdb.All()
	if err != nil {
		log.Fatalf("Failed to load sites. Error: %v", err)
	}
	for _, site := range sites {
		db := ws.NewEntryDB(site)
		dbc, err := db.Open()
		if err != nil {
			log.Fatalf("Failed to open DB. Error: %v", err)
		}
		err = dbc.Close()
		if err != nil {
			log.Fatalf("Failed to close DB. Error: %v", err)
		}
	}

	log.Fatal(http.ListenAndServe(conf.Main.WebEndpoint, ws.Router))
}
