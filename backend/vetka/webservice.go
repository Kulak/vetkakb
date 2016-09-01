package vetka

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"horse.lan.gnezdovi.com/vetkakb/backend/core"
)

// Helper functions

// WebSvc is a web service structure.
type WebSvc struct {
	Router  *httprouter.Router
	conf    *core.Configuration
	entryDB *core.EntryDB
	typeSvc *core.TypeService
}

// NewWebSvc creates new WebSvc structure.
func NewWebSvc(conf *core.Configuration, entryDB *core.EntryDB, typeSvc *core.TypeService) *WebSvc {

	ws := &WebSvc{
		Router:  httprouter.New(),
		conf:    conf,
		entryDB: entryDB,
		typeSvc: typeSvc,
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
	router.POST("/entry", ws.postEntry)
	router.GET("/api/recent", ws.getRecent)
	router.GET("/api/recent/:limit", ws.getRecent)
	router.GET("/api/entry/:entryID", ws.getFullEntry)
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
	// Post has one additional Parameter (EntryID) that will be set to zero.
	// If EntryID is zero, then it is a new item to be inserted.
	ws.handleWSEntryPost(w, r)
}

// postEntry updates existing entry and requires EntryID to exist.
func (ws WebSvc) postEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ws.handleWSEntryPost(w, r)
}

func (ws WebSvc) handleWSEntryPost(w http.ResponseWriter, r *http.Request) {
	var wse core.WSEntryPost
	err := ws.loadJSONBody(r, &wse)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	fmt.Printf("Got request with entry %v.\n", wse)
	//fmt.Printf("Request raw as string: %s\n", string(wse.Raw))
	var tp *core.TypeProvider
	tp, err = ws.typeSvc.Provider(wse.RawType)
	en := core.NewEntry(wse.EntryID, wse.Title, wse.Raw, wse.RawType)
	en.HTML, err = tp.ToHTML(wse.Raw)
	es := core.NewEntrySearch(wse.EntryID, wse.Tags)
	es.Plain, err = tp.ToPlain(wse.Raw)
	err = ws.entryDB.SaveEntry(en, es)
	if err != nil {
		ws.writeError(w, err.Error())
	} else {
		ws.writeJSON(w, en)
	}
}

func (ws WebSvc) getRecent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limitStr := p.ByName("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot parse limit.  Error: %v", err))
		return
	}
	if limit > 200 {
		limit = 200
	}
	entries, err := ws.entryDB.RecentHTMLEntries(limit)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to load recent HTML entries. Error: %v", err))
		return
	}
	ws.writeJSON(w, entries)
}

func (ws WebSvc) getFullEntry(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("entryID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot parse entryID.  Error: %v", err))
		return
	}
	entry, err := ws.entryDB.GetFullEntry(id)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get Entry with ID %v.  Error: %v", id, err))
		return
	}
	ws.writeJSON(w, entry)
}

// use POST to modify
