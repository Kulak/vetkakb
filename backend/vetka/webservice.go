package vetka

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"horse.lan.gnezdovi.com/vetkakb/backend/core"
	"horse.lan.gnezdovi.com/vetkakb/backend/edb"
	"horse.lan.gnezdovi.com/vetkakb/backend/sdb"
)

// Helper functions

// WebSvc is a web service structure.
type WebSvc struct {
	Router  *httprouter.Router
	conf    *core.Configuration
	siteDB  *sdb.SiteDB
	typeSvc *edb.TypeService
	store   *sessions.CookieStore
	// each key is a site's DBName.
	edbCache map[string]*edb.EntryDB
}

// NewWebSvc creates new WebSvc structure.
func NewWebSvc(conf *core.Configuration, siteDB *sdb.SiteDB, typeSvc *edb.TypeService) *WebSvc {

	gothic.Store = sessions.NewCookieStore([]byte("something-very-secret-blah-123!.;"))

	gplusKey := "214159873843-v6p3kmhikm62uc3j2paut5rsvkivod8v.apps.googleusercontent.com"
	gplusSecret := "0-eQESZIMdoKKn_2Xekl9e1b"
	goth.UseProviders(
		gplus.New(gplusKey, gplusSecret, fmt.Sprintf("%s/api/auth/callback?provider=gplus", conf.Main.SiteURL)),
	)

	ws := &WebSvc{
		Router:   httprouter.New(),
		conf:     conf,
		siteDB:   siteDB,
		typeSvc:  typeSvc,
		store:    sessions.NewCookieStore([]byte("moi-ochen-bolshoy-secret-123-!-21-13.")),
		edbCache: make(map[string]*edb.EntryDB),
	}

	// CRUD model in REST:
	//   create - POST
	//   read - GET
	//   update/replace - PUT
	//   update/modify - PATCH
	//   delete - DELETE

	router := ws.Router

	prefixes := []string{""}
	if len(ws.conf.Main.ClientPath) > 0 {
		prefixes = append(prefixes, "/client/:clientName")
	}
	for _, prefix := range prefixes {
		// serve static files
		router.GET(prefix+"/index.html", ws.siteHandler(ws.getIndex))
		router.GET(prefix+"/", ws.siteHandler(ws.getIndex))
		router.ServeFiles(prefix+"/vendors/*filepath", conf.WebDir("bower_components/"))
		router.ServeFiles(prefix+"/theme/*filepath", conf.WebDir("theme/"))
		// serve dynamic (site specific) content
		router.POST(prefix+"/binaryentry/", ws.siteHandler(ws.demandAdministrator(ws.postBinaryEntry)))
		router.GET(prefix+"/api/recent", ws.siteHandler(ws.getRecent))
		router.GET(prefix+"/api/recent/:limit", ws.siteHandler(ws.getRecent))
		router.GET(prefix+"/api/search/*query", ws.siteHandler(ws.getSearch))
		router.GET(prefix+"/api/entry/:entryID", ws.siteHandler(ws.getFullEntry))
		router.GET(prefix+"/api/rawtype/list", ws.siteHandler(ws.getRawTypeList))
		//router.HandlerFunc("GET", prefix+"/api/auth", ws.siteHandlerFunc(ws.beginAuthHandler(nil)))
		router.HandlerFunc("GET", prefix+"/api/auth", ws.siteHandlerFunc(gothic.BeginAuthHandler))
		router.GET(prefix+"/api/auth/callback", ws.getGplusCallback)
		// returns wsUserGet strucure usable for general web pages
		router.GET(prefix+"/api/session/user", ws.siteHandler(ws.wsUserGet))
		// for testing purpose of gothic cookie
		router.GET(prefix+"/api/session/gothic", ws.siteHandler(ws.demandAdministrator(ws.getGothicSession)))
		// for testing purpose of userId cookie
		router.GET(prefix+"/api/session/vetka", ws.siteHandler(ws.demandAdministrator(ws.getVetkaSession)))
		// allows to load RawTypeName "Binary/Image" as a link.
		router.GET(prefix+"/re/:entryID", ws.siteHandler(ws.getResourceEntry))
		// Enable access to source code files from web browser debugger
		router.ServeFiles(prefix+"/frontend/*filepath", http.Dir("frontend/"))
		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range r.Header {
				log.Println("Request HEADER key:", k, "value:", v)
			}
			msg := fmt.Sprintf("404 - File Not Found\n\nHost: %s\nURL: %s\nRequestURI: %s\n\n%v", r.Host, r.URL, r.RequestURI, r)
			ws.writeError(w, msg)
		})
	}
	return ws
}

// getGplusCallback is called by "Google Plus" OAuth2 API when user is authenticated.
// It creates new user if user is absent and sets "vetka" cookie with user id.
func (ws WebSvc) getGplusCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Processing google plus callback")
	gUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	siteIDStr := gothic.GetState(r)
	site, err := ws.siteDB.GetSiteByID(siteIDStr)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get site for SiteID %s.  Error: %s", siteIDStr, err))
		return
	}

	//log.Printf("Logged in user: %v", user)
	entryDB := ws.edbCache[site.DBName]
	user, err := entryDB.GetOrCreateUser(gUser)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}

	session, err := ws.store.Get(r, "vetka")
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to get vetka session store: %v", err))
		return
	}
	session.Values["userId"] = user.UserID
	session.Save(r, w)

	var fileName = fmt.Sprintf("http://%s/index.html", site.Host)
	if len(site.Path) > 0 {
		fileName = fmt.Sprintf("http://%s%s/%s/index.html", site.Host, ws.conf.Main.ClientPath, site.Path)
	}
	log.Printf("Redirect URL: %s", fileName)
	http.Redirect(w, r, fileName, 307)
}

// getSession is a study call to figure out what's inside gothic session
func (ws WebSvc) getGothicSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	providerName := "gplus"
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get provider: %v", err))
		return
	}

	session, err := gothic.Store.Get(r, gothic.SessionName)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get session: %v", err))
		return
	}

	if session.Values[gothic.SessionName] == nil {
		ws.writeError(w, "could not find a matching session for this request")
		return
	}

	sess, err := provider.UnmarshalSession(session.Values[gothic.SessionName].(string))
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot unmarshal session: %v", err))
		return
	}
	ws.writeJSON(w, sess)
	// Prints result like this:
	/*
		{"AuthURL":"https://accounts.google.com/o/oauth2/auth?access_type=offline\u0026client_id=214159873843-v6p3kmhikm62uc3j2paut5rsvkivod8v.apps.googleusercontent.com\u0026redirect_uri=http%3A%2F%2Fwww.gnezdovi.com%3A8080%2Fapi%2Fauth%2Fcallback%3Fprovider%3Dgplus\u0026response_type=code\u0026scope=profile+email+openid\u0026state=state","AccessToken":"","RefreshToken":"","ExpiresAt":"0001-01-01T00:00:00Z"}
	*/
}

func (ws WebSvc) getVetkaSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := ws.store.Get(r, "vetka")
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to get vetka session store: %v", err))
		return
	}
	ws.writeError(w, fmt.Sprintf("Vetka session userId: %v", session.Values["userId"]))
}

func (ws WebSvc) wsUserGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := ws.sessionWSUser(r)
	ws.writeJSON(w, u)
}

func (ws WebSvc) getRecent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limit, err := ws.getLimit(p)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	entries, err := entryDB.RecentHTMLEntries(limit)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to load recent HTML entries. Error: %v", err))
		return
	}
	ws.writeJSON(w, entries)
}

func (ws WebSvc) getIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ws.processTemplate(w, r, "index.html")
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
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	entries, err := entryDB.MatchEntries(query, limit)
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

func (ws WebSvc) getWSFullEntry(w http.ResponseWriter, r *http.Request, p httprouter.Params) *edb.WSFullEntry {
	idStr := p.ByName("entryID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot parse entryID.  Error: %v", err))
		return nil
	}
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	entry, err := entryDB.GetFullEntry(id)
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
	var wse edb.WSEntryPost
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
func (ws WebSvc) handleWSEntryPost(w http.ResponseWriter, r *http.Request, wse *edb.WSEntryPost, raw []byte) {
	var err error
	// we cannot log everything, because Raw may contain very large data
	fmt.Printf("Got request with entry id: %v, title: %v, rawTypeName: %v.\n", wse.EntryID, wse.Title, wse.RawTypeName)
	//fmt.Printf("Request raw as string: %s\n", string(wse.Raw))
	var tp *edb.TypeProvider
	tp, err = ws.typeSvc.ProviderByName(wse.RawTypeName)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	userID := ws.sessionUserID(r)
	en := edb.NewEntry(wse.EntryID, raw, tp.TypeNum, wse.RawContentType,
		wse.RawFileName, userID)
	en.HTML, err = tp.ToHTML(raw)
	es := edb.NewEntrySearch(wse.EntryID, wse.Title, wse.Tags)
	es.Plain, err = tp.ToPlain(raw)
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	err = entryDB.SaveEntry(en, es)
	if err != nil {
		ws.writeError(w, err.Error())
	} else {
		wen, err := entryDB.GetFullEntry(en.EntryID)
		if err != nil {
			ws.writeError(w, fmt.Sprintf("Cannot get Entry with ID %v.  Error: %v", en.EntryID, err))
			return
		}
		ws.writeJSON(w, wen)
	}
}
