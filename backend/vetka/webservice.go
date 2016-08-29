package vetka

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"horse.lan.gnezdovi.com/vetkakb/backend/core"
)

// Helper functions

// WebSvc is a web service structure.
type WebSvc struct {
	Router *httprouter.Router
	conf   *core.Configuration
}

// NewWebSvc creates new WebSvc structure.
func NewWebSvc(conf *core.Configuration) *WebSvc {

	ws := &WebSvc{
		Router: httprouter.New(),
		conf:   conf,
	}

	// CRUD model in REST:
	//   create - POST
	//   read - GET
	//   update/replace - PUT
	//   update/modify - PATCH
	//   delete - DELETE

	router := ws.Router
	router.GET("/index.html", ws.AddHeaders(http.FileServer(conf.WebDir("index.html"))))
	router.Handler("GET", "/", http.FileServer(conf.WebDir("/")))
	router.ServeFiles("/vendors/*filepath", conf.WebDir("bower_components/"))
	router.ServeFiles("/res/*filepath", conf.WebDir("res/"))
	router.PUT("/entry/", ws.putEntry)
	// Enable access to source code files from web browser debugger
	router.ServeFiles("/frontend/*filepath", http.Dir("frontend/"))

	return ws
}

// Handler functions

// AddHeaders adds custom HEADERs to index.html response using middleware style solution.
func (ws WebSvc) AddHeaders(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		/* Custom headers are easy to control here: */
		// fmt.Println("Adding header Access-Control-Allow-Credentials")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// w.Header().Set("Access-Control-Allow-Headers",
		// 	"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)
	}
}

// putEntry creates new entry and assigns it an EntryID.
func (ws WebSvc) putEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var wse WSEntryPut
	err := ws.loadJSONBody(r, &wse)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	fmt.Printf("Got request to create an entry with %v.\n", wse)
	fmt.Printf("Request raw as string: %s\n", string(wse.Raw))

	ws.writeError(w, "Can't save entries yet.")
}

// use POST to modify
