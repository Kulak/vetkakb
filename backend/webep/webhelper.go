package webep

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"horse.lan.gnezdovi.com/vetkakb/backend/edb"
	"horse.lan.gnezdovi.com/vetkakb/backend/sdb"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

func (ws WebSvc) getAuthProvider(site *sdb.Site) *gplus.Provider {
	// protocol could be inferred from request protocol
	siteURL := fmt.Sprintf("http://%s", site.Host)
	return gplus.New(gplusKey, gplusSecret, fmt.Sprintf("%s/api/auth/callback?provider=gplus", siteURL))
}

// Redirects client to authentication provider.
// State is not used.
func (ws WebSvc) beginAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Algorithm is based on :  gothic.BeginAuthHandler() call
	site := context.Get(r, "site").(*sdb.Site)
	provider := ws.getAuthProvider(site)
	state := gothic.SetState(r)

	sess, err := provider.BeginAuth(state)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to start authentication with provider: %s", err))
		return
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to get authentication url from provider: %s", err))
		return
	}

	session, _ := ws.store.Get(r, authSessionName)
	session.Values[authSessionName] = sess.Marshal()
	err = session.Save(r, w)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to save session in cookie store: %s", err))
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// getGplusCallback is called by "Google Plus" OAuth2 API when user is authenticated.
// It creates new user if user is absent and sets "vetka" cookie with user id.
func (ws WebSvc) getGplusCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	siteIDStr := gothic.GetState(r)
	site, err := ws.siteDB.GetSiteByID(siteIDStr)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get site for SiteID %s.  Error: %s", siteIDStr, err))
		return
	}

	// Algorithm is based on:     gothic.CompleteUserAuth(w, r)

	log.Println("Processing google plus callback")

	session, _ := ws.store.Get(r, authSessionName)

	if session.Values[authSessionName] == nil {
		ws.writeError(w, "could not find a matching session for this request")
		return
	}

	provider := ws.getAuthProvider(site)
	sess, err := provider.UnmarshalSession(session.Values[authSessionName].(string))
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to unmrarshall auth session: %s", err))
		return
	}

	_, err = sess.Authorize(provider, r.URL.Query())
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to authorize session url query: %s", err))
		return
	}

	// last line of gothic.CompleteUserAuth(w, r) algorithm
	gUser, err := provider.FetchUser(sess)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to fetch user from provider: %s", err))
		return
	}

	//log.Printf("Logged in user: %v", user)
	entryDB := ws.edbCache[site.DBName] // risky line; shall use context
	user, err := entryDB.GetOrCreateUser(gUser)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to get or create EntryDB user: %s", err))
		return
	}

	session, err = ws.store.Get(r, "vetka")
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to get vetka session store: %v", err))
		return
	}
	session.Values["userId"] = user.UserID
	err = session.Save(r, w)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to set vetka session store: %v", err))
		return
	}

	var fileName = fmt.Sprintf("http://%s/", site.Host)
	if len(site.Path) > 0 {
		fileName = fmt.Sprintf("http://%s%s/%s/", site.Host, ws.conf.Main.ClientPath, site.Path)
	}
	log.Printf("Redirect URL: %s", fileName)
	http.Redirect(w, r, fileName, 307)
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
		fmt.Printf("Failed to get user for userID %v from DB: %v", userID, err)
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

func (ws WebSvc) getLimit(limitStr string) (int64, error) {
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

func (ws WebSvc) getTime(timeStr string) (t time.Time, err error) {
	// Example of date in timeStr: "2016-09-22T15:05:22Z"
	t, err = time.ParseInLocation("2006-01-02T15:04:05Z", timeStr, time.UTC)
	return
}

// AddHeaders adds custom HEADERs to index.html response using middleware style solution.
func (ws WebSvc) addHeaders(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		/* Custom headers are easy to control here: */
		// fmt.Println("Adding header Access-Control-Allow-Credentials")
		// w.Header().Add("Access-Control-Allow-Origin", "http://webcache.googleusercontent.com")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// w.Header().Set("Access-Control-Allow-Headers",
		// 	"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)
	}
}

func (ws WebSvc) siteHandler(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := ws.setEdbContext(w, r)
		defer context.Clear(r)
		if err != nil {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	})
}

func (ws WebSvc) setEdbContext(w http.ResponseWriter, r *http.Request) (err error) {
	// log.Printf("Site request URL: %v; URL Path: %s, URL Host: %s, Host: %s",
	// 	r.URL, r.URL.Path, r.URL.Host, r.Host)
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
		ws.edbCache[site.DBName] = db
	}
	return db
}

// NewEntryDB creates new EntryDB based on web service context.
func (ws WebSvc) NewEntryDB(site *sdb.Site) *edb.EntryDB {
	db := edb.NewEntryDB(ws.conf.SQLDir("entrydb"), ws.conf.Main.DataRoot, site.DBName, ws.typeSvc)
	err := db.Upgrade()
	if err != nil {
		log.Fatalf("Failed to create or upgrade DB. Error: %v", err)
	}
	return db
}

// SiteProps provides basic interafce to template file.
type SiteProps struct {
	PageTitle string
	Theme     string
	GD        interface{}
}

// Builds a set of template files to be used in template execution.
// It uses theme as an entry point.
// It then loads site customization template.
//
// Global location is rooted in t-html/ directory tree.
// Site specific location rooted in dataDir/sitedb/t-html-s
//
// Global template must be present.  Site specific template might be optional.
// Global template is the 1st in returned result as it is used by template rendering sytem.
func (ws WebSvc) getWebTemplateFiles(r *http.Request, tFileName string) []string {
	result := []string{}
	// load theme specific global template file; template entry point
	site := context.Get(r, "site").(*sdb.Site)
	globalFile := ws.conf.TemplateThemeFile(site.Theme, tFileName)
	result = append(result, globalFile)
	// load site specifc customizations (file in t-html-s directory)
	siteFile, exists := site.WebTemplateFile(ws.conf.Main.DataRoot, tFileName)
	if exists {
		result = append(result, siteFile)
	} else {
		// load default site template included with theme
		defaultSiteFile := ws.conf.TemplateThemeFile(site.Theme, "default/"+tFileName)
		result = append(result, defaultSiteFile)
	}
	return result
}

func (ws WebSvc) processTemplate(w http.ResponseWriter, r *http.Request, tFileNames []string) {
	site := context.Get(r, "site").(*sdb.Site)
	log.Printf("Template file names: %s", tFileNames)
	// extracts name of the file with extension from the path
	baseName := path.Base(tFileNames[0])
	// New takes name of template (it needs to be name of parsed template base file name)
	// ParseFiles takes template file names to parse.  Abstract base template shall go 1st,
	t, err := template.New(baseName).ParseFiles(tFileNames...)
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

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}
