package vetka

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"horse.lan.gnezdovi.com/vetkakb/backend/edb"
	"horse.lan.gnezdovi.com/vetkakb/backend/sdb"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

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

func (ws WebSvc) demandAdministrator(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u := ws.sessionWSUser(r)
		if edb.Administrator.HasAccess(u.Clearances) {
			handler(w, r, p)
		} else {
			ws.writeError(w, http.StatusText(http.StatusUnauthorized))
		}
	}
}

func (ws WebSvc) sessionUserID(r *http.Request) (userID int64) {
	var err error
	var session *sessions.Session
	session, err = ws.store.Get(r, "vetka")
	if err != nil {
		fmt.Printf("Failed to get vetka session store: %v", err)
		return
	}
	userIDStr := session.Values["userId"]
	if userIDStr == nil {
		return 0
	}
	userID = userIDStr.(int64)
	return
}

// sessionUser returns current session user or guest if there is a problem.
// It always returns a valid user ID.
func (ws WebSvc) sessionWSUser(r *http.Request) (u *edb.WSUserGet) {
	var err error
	var userID int64
	userID = ws.sessionUserID(r)
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	u, err = entryDB.GetUser(userID)
	if err != nil {
		fmt.Printf("Failed to get user from DB: %v", err)
		u = edb.GuestWSUserGet
	}
	return
}

func (ws WebSvc) writeJSON(w http.ResponseWriter, v interface{}) {
	encoded, err := json.Marshal(v)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to encode to JSON: %v", err))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", encoded)
	}
}

func (ws WebSvc) writeError(w http.ResponseWriter, msg string) {
	log.Printf("500: %v", msg)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(500)
	fmt.Fprint(w, msg)
}

func (ws WebSvc) writeResult(w http.ResponseWriter, v interface{}, err error, msg string) {
	if err != nil {
		ws.writeError(w, fmt.Sprintf("%s: %v", msg, err))
	} else {
		ws.writeJSON(w, v)
	}
}

func (ws WebSvc) loadJSONBody(rBody io.ReadCloser, v interface{}) error {
	// load post data
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(rBody)
	if err != nil {
		return fmt.Errorf("Failed to read request body: %v", err)
	}
	fmt.Printf("request body: %s\n", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("Failed to unmarshall request body: %v", err)
	}
	return nil
}

func (ws WebSvc) getLimit(p httprouter.Params) (int64, error) {
	limitStr := p.ByName("limit")
	if limitStr == "" {
		limitStr = "30"
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Cannot parse limit.  Error: %v", err)
	}
	if limit > 200 {
		limit = 200
	}
	return limit, nil
}

func (ws WebSvc) siteHandler(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := ws.setEdbContext(w, r)
		defer context.Clear(r)
		if err != nil {
			return
		}
		handler(w, r, p)
	}
}

func (ws WebSvc) siteHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ws.setEdbContext(w, r)
		defer context.Clear(r)
		if err != nil {
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (ws WebSvc) setEdbContext(w http.ResponseWriter, r *http.Request) (err error) {
	log.Printf("Site request URL: %v; URL Path: %s, URL Host: %s, Host: %s",
		r.URL, r.URL.Path, r.URL.Host, r.Host)
	var path = ""
	var site *sdb.Site
	if strings.HasPrefix(r.URL.Path, ws.conf.Main.ClientPath) {
		// if ws.conf.Main.ClientPath is set to /client, then valid URL would be
		// 	/client/doha/api/recent
		// In this case array split is:
		// [0] is ""
		// [1] is client
		// [2] is client name that we need to path as path
		paths := strings.Split(r.URL.Path, "/")
		if len(paths) > 2 {
			path = paths[2]
		}
		site, err = ws.siteDB.GetSite(r.Host, path)
	} else {
		site, err = ws.siteDB.GetSite(r.Host, path)
		if site.SiteID == 0 {
			// try to load default site
			site, err = ws.siteDB.GetSite("", "")
		}
	}
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot locate site based on host %s and path %s.", r.URL.Host, path))
		return
	}
	if site.Path != "" {
		site.ZonePath = ws.conf.Main.ClientPath + "/" + site.Path
	}
	db := ws.getEdb(site)
	context.Set(r, "edb", db)
	context.Set(r, "site", site)
	return
}

// getEdb returns EntryDB for site.
// It uses edbCache or populates it.
func (ws WebSvc) getEdb(site *sdb.Site) *edb.EntryDB {
	db, ok := ws.edbCache[site.DBName]
	if !ok {
		log.Printf("Caching EntityDB %s", site.DBName)
		db = ws.NewEntryDB(site)
		err := db.Open()
		if err != nil {
			log.Printf("Cannot load EntryDB for site: %v", err)
			return nil
		}
		ws.edbCache[site.DBName] = db
	}
	return db
}

// NewEntryDB creates new EntryDB based on web service context.
func (ws WebSvc) NewEntryDB(site *sdb.Site) *edb.EntryDB {
	return edb.NewEntryDB(ws.conf.SQLDir("entrydb"), ws.conf.Main.DataRoot,
		site.DBName, ws.typeSvc)
}

// SiteProps provides basic interafce to template file.
type SiteProps struct {
	PageTitle string
	Theme     string
	GD        interface{}
}

func (ws WebSvc) processTemplate(w http.ResponseWriter, r *http.Request, tFileName string) {
	site := context.Get(r, "site").(*sdb.Site)
	indexTFile := ws.conf.TemplateThemeFile(site.Theme, tFileName)
	// New takes name of template (can be anything)
	// ParseFiles takes template file names to parse.  Abstract base template shall go 1st,
	t, err := template.New(tFileName).ParseFiles(indexTFile)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot parse template.  Error: %s", err))
		return
	}
	err = t.Execute(w, site)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot execute template.  Error: %s", err))
		return
	}
}
