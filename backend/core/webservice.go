package core

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Helper functions

type WebSvc struct {
	Router *httprouter.Router
	conf   *Configuration
}

func NewWebSvc(conf *Configuration) *WebSvc {

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
	router.GET("/index.html", ws.AddHeaders(http.FileServer(conf.WebDir("www/index.html"))))
	router.Handler("GET", "/", http.FileServer(conf.WebDir("www/")))
	router.ServeFiles("/vendors/*filepath", conf.WebDir("www/bower_components/"))

	return ws
}

// Handler functions

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
