package vetka

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	router.PUT("/binaryentry/", ws.putBinaryEntry)
	router.POST("/entry", ws.postEntry)
	router.GET("/api/recent", ws.getRecent)
	router.GET("/api/recent/:limit", ws.getRecent)
	router.GET("/api/search/*query", ws.getSearch)
	router.GET("/api/entry/:entryID", ws.getFullEntry)
	router.GET("/api/rawtype/list", ws.getRawTypeList)
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
	// ws.handleWSEntryPost(w, r)
	ws.writeError(w, "PUT entry is currently not implemented.")
}

// postEntry updates existing entry and requires EntryID to exist.
func (ws WebSvc) postEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//ws.handleWSEntryPost(w, r)
	ws.writeError(w, "POST entry is currently not implemented.")
}

func (ws WebSvc) getRecent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limit, err := ws.getLimit(p)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	entries, err := ws.entryDB.RecentHTMLEntries(limit)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to load recent HTML entries. Error: %v", err))
		return
	}
	ws.writeJSON(w, entries)
}

// getMatch searches for entries matching criteria.
func (ws WebSvc) getSearch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limit, err := ws.getLimit(p)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	query := p.ByName("query")
	if len(query) < 2 {
		ws.writeError(w, fmt.Sprintf("Query is not supplied (len:%v).", len(query)))
		return
	}
	query = query[1:]
	entries, err := ws.entryDB.MatchEntries(query, limit)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to match HTML entries. Error: %v", err))
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

func (ws WebSvc) getRawTypeList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type WSRawType struct {
		TypeNum int
		Name    string
	}
	list := []WSRawType{}
	for k, v := range ws.typeSvc.List() {
		list = append(list, WSRawType{TypeNum: k, Name: v.Name})
	}
	ws.writeJSON(w, list)
}

func (ws WebSvc) putBinaryEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("receiving binary data")
	ws.handleAnyWSEntryPost(w, r)
}

func (ws WebSvc) handleAnyWSEntryPost(w http.ResponseWriter, r *http.Request) {
	// Standard multi-part PULL or POST consists of 2 parts:
	// Part 1 is a JSON message
	// Part 2 is a binary message representing Raw column value of Entry table.

	mr, err := r.MultipartReader()
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Error reading multipart header: %v", err))
		return
	}

	// the goal of this loop is to populate wse variable
	// with JSON and RAW data.
	var wse core.WSEntryPost

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			log.Printf("Reached EOF in multi-part read.")
			break
		} else if err != nil {
			ws.writeError(w, fmt.Sprintf("Error reading multi-part message: %v", err))
			return
		}

		// Here is a typical output for log line below.
		// log.Printf("Header: %v", p.Header)
		// Header for entry: map[Content-Disposition:[form-data; name="entry"]]
		// Header for image: Header: map[Content-Disposition:[form-data; name="rawFile"; filename="1upatime-pronoun-icon.png"] Content-Type:[image/png]]

		// FormName on javaScript side corresponds to 1st argument of FormData.append
		log.Printf("Part file name: %s, form name: %s\n", p.FileName(), p.FormName())
		switch p.FormName() {
		case "entry":
			// decode standard JSON message: {"title":"","raw":null,"rawType":4,"tags":""}
			err := ws.loadJSONBody(p, &wse)
			if err != nil {
				ws.writeError(w, fmt.Sprintf("Error reading entry part: %v", err))
				return
			}
		case "rawFile":
			// read raw bytes
			// var origFileName = p.FileName()
			content, err := ioutil.ReadAll(p)
			if err != nil {
				ws.writeError(w, fmt.Sprintf("Error reading rawFile part: %v", err))
				return
			}
			// Raw assignment is our primary goal
			wse.Raw = content
			// Write a temporary diagnostics file
			targetFile := ws.conf.DataFile("last-uploaded.jpg")
			err = ioutil.WriteFile(targetFile, content, 0777)
			if err != nil {
				// don't write error message, because this is a diagnostics code; just log
				log.Printf("Failed to save receipt image in the database: %v", err)
			}
		default:
			log.Printf("unrecognized FormName: %v", p.FormName())
		}
	} // end of loop for multiple parts
	// validate that we received expected data
	if wse.RawType == 0 {
		ws.writeError(w, "RawType is not received.")
		return
	}
	if wse.Raw == nil {
		ws.writeError(w, "Raw payload is not received.")
		return
	}
	ws.handleWSEntryPost(w, r, &wse)
}

// handleWSEntryPost inserts or updates Entry using standard algorithm.
func (ws WebSvc) handleWSEntryPost(w http.ResponseWriter, r *http.Request, wse *core.WSEntryPost) {
	var err error
	// we cannot log everything, because Raw may contain very large data
	fmt.Printf("Got request with entry id: %v, title: %v, rawType: %v.\n", wse.EntryID, wse.Title, wse.RawType)
	//fmt.Printf("Request raw as string: %s\n", string(wse.Raw))
	var tp *core.TypeProvider
	tp, err = ws.typeSvc.Provider(wse.RawType)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
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
