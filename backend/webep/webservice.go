package webep

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
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

var authSessionName = "vetka_auth"
var gplusKey = "214159873843-v6p3kmhikm62uc3j2paut5rsvkivod8v.apps.googleusercontent.com"
var gplusSecret = "0-eQESZIMdoKKn_2Xekl9e1b"

// NewWebSvc creates new WebSvc structure.
func NewWebSvc(conf *core.Configuration, siteDB *sdb.SiteDB, typeSvc *edb.TypeService) *WebSvc {

	gothic.Store = sessions.NewCookieStore([]byte("something-very-secret-blah-123!.;"))

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

	router.GET("/robots.txt", ws.getRobots)
	prefixes := []string{""}
	if len(ws.conf.Main.ClientPath) > 0 {
		// ClientPath form is "/cl" or "/client"
		prefixes = append(prefixes, fmt.Sprintf("%s/:clientName", ws.conf.Main.ClientPath))
	}
	log.Println("Registering URL prefixes: ", prefixes)
	for _, prefix := range prefixes {
		// serve static files
		router.GET(prefix+"/index.html", ws.siteHandler(ws.getIndex))
		router.GET(prefix+"/", ws.siteHandler(ws.getIndex))
		router.GET(prefix+"/app/*ignoredPageName", ws.siteHandler(ws.getIndex))
		// load by slug
		router.GET(prefix+"/s/*ignoredPageName", ws.siteHandler(ws.getIndex))
		//router.GET(prefix+"/app/e/:entryID/*ignoredSlug", ws.siteHandler(ws.getIndex))
		router.ServeFiles(prefix+"/theme/*filepath", conf.WebDir("theme/"))
		// serve dynamic (site specific) content
		router.POST(prefix+"/binaryentry/", ws.siteHandler(ws.demandAdministrator(ws.postBinaryEntry)))
		// generates sitemap XML
		router.GET(prefix+"/api/sitemap", ws.siteHandler(ws.generateSitemap))
		router.GET(prefix+"/api/recent/:limit", ws.siteHandler(ws.getRecent))
		router.GET(prefix+"/api/recent/:limit/:end", ws.siteHandler(ws.getRecent))
		router.GET(prefix+"/api/search/*query", ws.siteHandler(ws.getSearch))
		router.GET(prefix+"/api/entry/:entryID", ws.siteHandler(ws.getFullEntry))
		router.GET(prefix+"/api/s/:slug", ws.siteHandler(ws.getFullEntryBySlug))
		router.GET(prefix+"/api/rawtype/list", ws.siteHandler(ws.getRawTypeList))
		// site can be extracted when starting authentication
		router.HandlerFunc("GET", prefix+"/api/auth", ws.siteHandlerFunc(ws.beginAuthHandler))
		// callback cannot maintian zones, so site has to be loaded from state
		router.GET(prefix+"/api/auth/callback", ws.getGplusCallback)
		// returns wsUserGet strucure usable for general web pages
		router.GET(prefix+"/api/session/user", ws.siteHandler(ws.wsUserGet))
		// for testing purpose of gothic cookie
		router.GET(prefix+"/api/session/gothic", ws.siteHandler(ws.demandAdministrator(ws.getGothicSession)))
		// for testing purpose of userId cookie
		router.GET(prefix+"/api/session/vetka", ws.siteHandler(ws.demandAdministrator(ws.getVetkaSession)))
		// to get a quick public list of currently registered users; not really for display
		router.GET(prefix+"/api/users", ws.siteHandler(ws.demandAdministrator(ws.getUsers)))
		// allows to load RawTypeName "Binary/Image" as a link.
		router.GET(prefix+"/re/:entryID", ws.siteHandler(ws.getResourceEntry))
		// serve files under site's www/res directory
		router.GET(prefix+"/res/*filepath", ws.siteHandler(ws.serveResFile))
		// serve site files
		router.GET(prefix+"/sitemaps/*filepath", ws.siteHandler(ws.serveSitemapsFile))
		// Enable access to source code files from web browser debugger
		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// for k, v := range r.Header {
			// 	log.Println("Request HEADER key:", k, "value:", v)
			// }
			// log.Println("RquestURI:", r.RequestURI)
			msg := fmt.Sprintf("404 - File Not Found\n\nHost: %s\nURL: %s", r.Host, r.URL)
			ws.writeError(w, msg)
		})
	}
	router.ServeFiles("/vendors/*filepath", conf.WebDir("bower_components/"))
	router.ServeFiles("/frontend/*filepath", http.Dir("frontend/"))
	// site specific URLs
	sites, err := ws.siteDB.All()
	if err != nil {
		log.Fatalf("Failed to load sites: %v", err)
	}
	for _, site := range sites {
		if site.Path != "" {
			// redirect is only implemented for domain hosting, not zone level
			log.Printf("Redirect configuration is skipping site with %s path, host: %v, siteID: %v",
				site.Path, site.Host, site.SiteID)
			continue
		}
		edb := ws.getEdb(site)
		paths, err := edb.GetUniqueRedirectPaths()
		if err != nil {
			log.Fatalf("Failed to get unique redirect paths for site %v.  Error: %s", site.SiteID, err)
		}
		for _, path := range paths {
			// we can only redirect absolute path at domain level
			if strings.HasPrefix(path, "/") {
				log.Printf("Redirect for %s path", path)
				router.GET(path+"/*filepath", ws.siteHandler(ws.getRedirect))
				continue
			}
			log.Printf("Domain redirect path %s does not start with /", path)
		}
	}
	return ws
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

func (ws WebSvc) generateSitemap(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	entries, err := entryDB.AllHTMLEntries()
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}

	site := context.Get(r, "site").(*sdb.Site)

	sm := stm.NewSitemap()
	// controls URL generated in sitemap file
	sm.SetDefaultHost(fmt.Sprintf("http://%s", site.Host))
	// controls file system root location of the generated sitemap
	// it appends "sitemaps/sitemap.xml" and "sitemaps/sitemap1.xml"
	sm.SetPublicPath(site.WebFile(ws.conf.Main.DataRoot, ""))
	sm.SetCompress(false)
	sm.Create()
	m := make(map[string]string)
	for _, entry := range entries {
		path := fmt.Sprintf("%s/s/%s", site.ZonePath, entry.Slug)
		if entry.Slug == "" {
			m[strconv.Itoa(int(entry.EntryID))] = fmt.Sprintf("Entry has empty slug.  Title: %s", entry.Title)
			continue
		}
		sm.Add(stm.URL{"loc": path, "changefreq": "monthly", "mobile": true, "priority": 0.5})
	}
	sm.Finalize().PingSearchEngines()

	m["site.Host"] = site.Host
	ws.writeJSON(w, m)
}

func (ws WebSvc) getRecent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limitStr := p.ByName("limit")
	limit, err := ws.getLimit(limitStr)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	endStr := p.ByName("end")
	end := time.Now().Add(24 * time.Hour)
	if endStr != "" {
		end, err = ws.getTime(endStr)
		if err != nil {
			ws.writeError(w, err.Error())
			return
		}
	}

	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	entries, err := entryDB.RecentHTMLEntries(limit, end)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Failed to load recent HTML entries. Error: %v", err))
		return
	}
	ws.writeJSON(w, entries)
}

func (ws WebSvc) getIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fileNames := ws.getWebTemplateFiles(r, "index.html")
	ws.processTemplate(w, r, fileNames)
}

func (ws WebSvc) getRobots(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, ws.conf.WebFile("robots.txt"))
}

// getMatch searches for entries matching criteria.
func (ws WebSvc) getSearch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	limitStr := p.ByName("limit")
	limit, err := ws.getLimit(limitStr)
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

func (ws WebSvc) getFullEntryBySlug(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("slug")
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	entry, err := entryDB.GetFullEntryBySlug(slug)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get Entry with slug %v.  Error: %v", slug, err))
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
	fmt.Printf("Got request with entry id: %v, title: %v, rawTypeName: %v, titleIcon: %s, slug: %s.\n",
		wse.EntryID, wse.Title, wse.RawTypeName, wse.TitleIcon, wse.Slug)
	//fmt.Printf("Request raw as string: %s\n", string(wse.Raw))
	var tp *edb.TypeProvider
	tp, err = ws.typeSvc.ProviderByName(wse.RawTypeName)
	if err != nil {
		ws.writeError(w, err.Error())
		return
	}
	userID := ws.sessionUserID(r)
	en := edb.NewEntry(wse.EntryID, raw, tp.TypeNum, wse.RawContentType,
		wse.RawFileName, wse.TitleIcon, wse.Intro, wse.Slug, userID)
	en.HTML, err = tp.ToHTML(raw)
	es := edb.NewEntrySearch(wse.EntryID, wse.Title, wse.Tags)
	es.Plain, err = tp.ToPlain(raw)
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	var redirect *edb.Redirect
	err = entryDB.SaveEntry(en, es, redirect)
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

func (ws WebSvc) getRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	site := context.Get(r, "site").(*sdb.Site)
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	path := r.URL.Path
	entryID, err := entryDB.GetRedirectEntryID(path)
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot find redirect ID for path %s", path))
		return
	}
	to := fmt.Sprintf("http://%s/api/entry/%v", site.Host, entryID)
	log.Printf("Redirecting requested %s to %s", path, to)
	http.Redirect(w, r, to, 301)
}

func (ws WebSvc) getUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	entryDB := context.Get(r, "edb").(*edb.EntryDB)
	users, err := entryDB.GetUsers()
	if err != nil {
		ws.writeError(w, fmt.Sprintf("Cannot get users list: %s", err))
		return
	}
	ws.writeJSON(w, users)
}

func (ws WebSvc) serveResFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ws.serveSiteWebSubDirFile(w, r, p, "res")
}

func (ws WebSvc) serveSitemapsFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ws.serveSiteWebSubDirFile(w, r, p, "sitemaps")
}

func (ws WebSvc) serveSiteWebSubDirFile(w http.ResponseWriter, r *http.Request, p httprouter.Params, subdir string) {
	site := context.Get(r, "site").(*sdb.Site)
	fp := p.ByName("filepath")
	fn := filepath.Join(subdir, fp)
	sfn := site.WebFile(ws.conf.Main.DataRoot, fn)
	http.ServeFile(w, r, sfn)
}
