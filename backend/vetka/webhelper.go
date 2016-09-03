package vetka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

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

func (ws WebSvc) loadJSONBody(r *http.Request, v interface{}) error {
	// load post data
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(r.Body)
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
