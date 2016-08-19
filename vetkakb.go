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
	"net/http"
	"os/user"
	"strconv"
	"syscall"

	"horse.lan.gnezdovi.com/vetkakb/core"

	"github.com/sevlyar/go-daemon"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "cf", "", "Configuration file name")
	flag.Parse()
	log.Println("*** Starting rashodi ***")
	log.Printf("Flag -cf %v", configFile)

	conf, err := core.LoadConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	context := &daemon.Context{
		PidFileName: conf.Main.PidFileName,
		PidFilePerm: 0644,
		LogFileName: conf.Main.LogFileName,
		LogFilePerm: 0640,
		Credential:  &syscall.Credential{},
		//WorkDir:     "/",
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

	// child, err := context.Reborn()
	// if err != nil {
	// 	log.Fatalf("Failed to fork rashodi.  Error: %v", err)
	// }
	// if child != nil {
	// 	// parent does nothing with its child
	// 	log.Println("Master process ended.")
	// 	return
	// }

	// child code beyond this point
	log.Println("Child process continues.")
	defer context.Release()

	err = conf.InitializeFilesystem()
	if err != nil {
		log.Fatalf("Failed to InitializeFilesystem.  Error: %v", err)
	}

	log.Println("Startign web service")
	ws := core.NewWebSvc(conf)
	log.Fatal(http.ListenAndServe(conf.Main.WebEndpoint, ws.Router))
}
