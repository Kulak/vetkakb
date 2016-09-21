package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	"horse.lan.gnezdovi.com/vetkakb/backend/edb"

	_ "github.com/mattn/go-sqlite3"
)

// RslpArticle contains RSLP article info.
type RslpArticle struct {
	ArticleID  int64
	Path       string
	CategoryID int64
	Title      string
	Created    time.Time
	Updated    *time.Time
	Published  *time.Time
	Intro      string
	IconURL    string

	Content *string
	Tags    *string
}

var ownerID int64 = 1
var mp *edb.TypeProvider
var ed *edb.EntryDB

// Example:
// 	go run rslp-imp/rslp-imp.go -n default
func main() {
	log.Println("rslp-import into vetka")

	var dataDir string
	var dbName string
	flag.StringVar(&dbName, "n", "www.rebeccaslp.com", "entry database file name and directory")
	flag.StringVar(&dataDir, "d", "data", "data directory")
	flag.Parse()

	// globals
	mp = edb.MarkdownProvider()
	ed = edb.NewEntryDB("no-sql", dataDir, dbName, nil)
	err := ed.Open()
	if err != nil {
		log.Fatalf("Failed to open entry db: %s", err)
	}
	defer ed.Close()

	rslp, err := sql.Open("sqlite3", "data/rslp.sqlite")
	if err != nil {
		log.Fatalf("Failed to open 'data/rslp.sqlite' database. Error: %v", err)
	}
	defer rslp.Close()

	query := `
	SELECT
		a.articleId, a.path, a.categoryId, a.title,
		a.created, a.updated, a.published, a.intro,
		a.iconUrl, c.Content, c.Tags
	from articleInfo a
	inner join article c on a.ArticleID = c.ArticleID
	order by a.ArticleID asc
	`
	var rows *sql.Rows
	rows, err = rslp.Query(query)
	if err != nil {
		log.Fatalf("Failed to load rslp records: %s", err)
	}
	var r = RslpArticle{}
	for rows.Next() {
		err = rows.Scan(
			&r.ArticleID, &r.Path, &r.CategoryID, &r.Title,
			&r.Created, &r.Updated, &r.Published, &r.Intro,
			&r.IconURL, &r.Content, &r.Tags,
		)
		if err != nil {
			log.Fatalf("Failed to scan: %s", err)
		}
		insertRecord(&r)
		// // resolve number to a name
		// r.RawTypeName = edb.rawTypes.NameByNum(rawType)
		// result = append(result, r)
	}
	log.Println("All Done")
}

func insertRecord(r *RslpArticle) {
	log.Println("Record: ", r.Created, r.Updated, r.Published)

	var err error
	raw := []byte{}
	if r.Content != nil {
		raw = []byte(*r.Content)
	}

	en := edb.Entry{}
	// in alphabetical order
	en.Created = r.Created
	en.EntryID = 0
	en.HTML, err = mp.ToHTML(raw)
	if err != nil {
		log.Fatalf("Error converting to HTML: %s", err)
	}
	en.Intro = r.Intro
	en.OwnerFK = ownerID
	en.Published = r.Published
	en.Raw = raw
	en.RawContentType = ""   // not used in Markdown
	en.RawFileName = ""      // not used in Markdown
	en.RawType = 3           // 3 is markdown type
	en.RequiredClearance = 8 // 8 is administrator mask
	en.Updated = r.Created
	en.TitleIcon = r.IconURL
	if r.Updated != nil {
		en.Updated = *r.Updated
	}

	es := edb.EntrySearch{}
	es.EntryFK = 0
	es.Plain, err = mp.ToPlain(raw)
	if err != nil {
		log.Fatalf("Error converting to Plain: %s", err)
	}
	es.Tags = ""
	if r.Tags != nil {
		es.Tags = *r.Tags
	}
	es.Title = r.Title

	rr := edb.Redirect{}
	rr.Path = r.Path
	rr.StatusCode = 0

	err = ed.SaveEntry(&en, &es, &rr)
	if err != nil {
		log.Fatalf("Failed to insert entry: %s", err)
	}
}
