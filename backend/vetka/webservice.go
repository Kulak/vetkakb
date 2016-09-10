package vetka

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
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

	gothic.Store = sessions.NewFilesystemStore("/tmp/vetka-sessions/", []byte("goth-example"))

	gplusKey := "214159873843-v6p3kmhikm62uc3j2paut5rsvkivod8v.apps.googleusercontent.com"
	gplusSecret := "0-eQESZIMdoKKn_2Xekl9e1b"
	goth.UseProviders(
		gplus.New(gplusKey, gplusSecret, fmt.Sprintf("%s/api/auth/callback?provider=gplus", conf.Main.SiteURL)),
	)

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
	router.POST("/binaryentry/", ws.postBinaryEntry)
	router.GET("/api/recent", ws.getRecent)
	router.GET("/api/recent/:limit", ws.getRecent)
	router.GET("/api/search/*query", ws.getSearch)
	router.GET("/api/entry/:entryID", ws.getFullEntry)
	router.GET("/api/rawtype/list", ws.getRawTypeList)
	router.HandlerFunc("GET", "/api/auth", gothic.BeginAuthHandler)
	router.GET("/api/auth/callback", ws.getGplusCallback)
	// allows to load RawTypeName "Binary/Image" as a link.
	router.GET("/re/:entryID", ws.getResourceEntry)
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

func (ws WebSvc) getGplusCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	log.Printf("Logged in user: %v", user)
	fileName := ws.conf.Main.SiteURL + "/index.html"
	http.Redirect(w, r, fileName, 307)
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
	entry := ws.getWSFullEntry(w, r, p)
	if entry == nil {
		return
	}
	ws.writeJSON(w, entry)
}

func (ws WebSvc) getResourceEntry(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	entry := ws.getWSFullEntry(w, r, p)
	if entry == nil {
		return
	}
	if entry.RawTypeName == "Binary/Image" {
		w.Header().Set("Content-Type", "image/png")
		w.Write(entry.Raw)
		return
	}
	ws.writeError(w, "re/:entryId url path represents only binary resource")
}

func (ws WebSvc) getWSFullEntry(w http.ResponseWriter, r *http.Request, p httprouter.Params) *core.WSFullEntry {
	idStr := p.ByName("entryID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot parse entryID.  Error: %v", err))
		return nil
	}
	entry, err := ws.entryDB.GetFullEntry(id)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get Entry with ID %v.  Error: %v", id, err))
		return nil
	}
	return entry
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

func (ws WebSvc) postBinaryEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	var raw []byte
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
		log.Printf("Part form name: %s, file: %s, content-type: %s\n", p.FormName(), p.FileName(), p.Header.Get("Content-Type"))
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
			// Raw assignment is our primary goal
			raw, err = ioutil.ReadAll(p)
			if err != nil {
				ws.writeError(w, fmt.Sprintf("Error reading rawFile part: %v", err))
				return
			}
			wse.RawContentType = p.Header.Get("Content-Type")
			wse.RawFileName = p.FileName()
			// Write a temporary diagnostics file
			targetFile := ws.conf.DataFile("last-uploaded.jpg")
			err = ioutil.WriteFile(targetFile, raw, 0777)
			if err != nil {
				// don't write error message, because this is a diagnostics code; just log
				log.Printf("Failed to save receipt image in the database: %v", err)
			}
		default:
			log.Printf("unrecognized FormName: %v", p.FormName())
		}
	} // end of loop for multiple parts
	// validate that we received expected data
	if wse.RawTypeName == "" {
		ws.writeError(w, "RawTypeName is not received.")
		return
	}
	if raw == nil {
		ws.writeError(w, "Raw payload is not received.")
		return
	}
	ws.handleWSEntryPost(w, r, &wse, raw)
}

// handleWSEntryPost inserts or updates Entry using standard algorithm.
func (ws WebSvc) handleWSEntryPost(w http.ResponseWriter, r *http.Request, wse *core.WSEntryPost, raw []byte) {
	var err error
	// we cannot log everything, because Raw may contain very large data
	fmt.Printf("Got request with entry id: %v, title: %v, rawTypeName: %v.\n", wse.EntryID, wse.Title, wse.RawTypeName)
	//fmt.Printf("Request raw as string: %s\n", string(wse.Raw))
	var tp *core.TypeProvider
	tp, err = ws.typeSvc.ProviderByName(wse.RawTypeName)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	en := core.NewEntry(wse.EntryID, raw, tp.TypeNum, wse.RawContentType, wse.RawFileName)
	en.HTML, err = tp.ToHTML(raw)
	es := core.NewEntrySearch(wse.EntryID, wse.Title, wse.Tags)
	es.Plain, err = tp.ToPlain(raw)
	err = ws.entryDB.SaveEntry(en, es)
	if err != nil {
		ws.writeError(w, err.Error())
	} else {
		wen, err := ws.entryDB.GetFullEntry(en.EntryID)
		if err != nil {
			ws.writeError(w, fmt.Sprintf("Cannot get Entry with ID %v.  Error: %v", en.EntryID, err))
			return
		}
		ws.writeJSON(w, wen)
	}
}
