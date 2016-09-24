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
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
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
	var signal string
	flag.StringVar(&configFile, "cf", "/usr/local/etc/vetkakb.ini", "Configuration file name")
	flag.BoolVar(&consoleMode, "c", false, "Console mode (non-daemon)")
	flag.StringVar(&signal, "s", "", `send signal to daemon
			quit - graceful shutdown`)
	flag.Parse()
	log.Println("*** Starting vetkakb ***")
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
		deamonRun(conf, signal)
	}
}

func deamonRun(conf *core.Configuration, signal string) {
	var err error

	daemon.AddCommand(daemon.StringFlag(&signal, "quit"), syscall.SIGQUIT, quitHandler)

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

	if len(daemon.ActiveFlags()) > 0 {
		// process signals
		d, err := context.Search()
		if err != nil {
			fmt.Println("Unable send signal to vetkakb: ", err)
			return
		}
		err = daemon.SendCommands(d)
		if err != nil {
			fmt.Println("Failed to send command to vetkakb: ", err)
		} else {
			fmt.Println("Sent signal to hwpushd.")
		}
		return
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
		log.Fatalf("Failed to fork vetkakb.  Error: %v", err)
	}
	if child != nil {
		// parent does nothing with its child
		log.Println("Master process ended.")
		return
	}

	// child code beyond this point
	log.Println("Child process continues.")
	defer context.Release()

	go func() {
		childRun(conf)
	}()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Signal handler failed:", err)
	}
	log.Println("- - - vetka quit - - -")
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

// var (
// 	stop = make(chan struct{})
// 	done = make(chan struct{})
// )

// Graceful shutdown
func quitHandler(sig os.Signal) error {
	log.Println("quitting gracefully...")
	// stop <- struct{}{}
	// <-done
	return nil // daemon.ErrStop
}
