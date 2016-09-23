package edb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kulak/sqlitemaint"
	"github.com/markbates/goth"
)

// EntryDB is a service to interact with EntryDB database.
type EntryDB struct {
	db       *sql.DB
	sqlDir   string
	dataDir  string
	dbDir    string
	dbName   string
	rawTypes *TypeService
}

// NewEntryDB creates DB service, but does not open connection.
func NewEntryDB(sqlDir, dataDir, dbName string, rawTypes *TypeService) *EntryDB {
	return &EntryDB{
		sqlDir:   sqlDir,
		dataDir:  dataDir,
		dbDir:    filepath.Join(dataDir, dbName),
		dbName:   dbName,
		rawTypes: rawTypes,
	}
}

// Open opens DB connection.
func (edb *EntryDB) Open() error {
	if edb.db != nil {
		return nil
	}
	err := os.MkdirAll(edb.dbDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create a data directory %s due to error: %v", edb.dbDir, err)
	}
	dbFileName := filepath.Join(edb.dbDir, edb.dbName+".db")
	log.Printf("Entry DB file name: %s", dbFileName)
	_, err = sqlitemaint.UpgradeSQLite(dbFileName, edb.sqlDir)
	if err != nil {
		return fmt.Errorf("Failed to upgrade entry DB %s.  Error: %v", dbFileName, err)
	}
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return fmt.Errorf("Failed to open database %v. Error: %v", dbFileName, err)
	}
	edb.db = db
	return nil
}

// Close closes DB connection if it is open.
func (edb *EntryDB) Close() {
	if edb.db != nil {
		edb.db.Close()
		edb.db = nil
	}
}

// SaveEntry inserts new entres to DB and updates existing.
// If en.EntryID is zero, then it is considered new.
// If en.EntryID is not zero, then it is updated.
// When new objects are inserted en.EntryID and es.EntryFK
// are set to inserted ID value.
func (edb *EntryDB) SaveEntry(en *Entry, es *EntrySearch, r *Redirect,
) (err error) {
	if edb.db == nil {
		return fmt.Errorf("Database connection is closed.")
	}
	if en.EntryID != es.EntryFK {
		return fmt.Errorf("EntryID %v does not match EntryFK %v", en.EntryID, es.EntryFK)
	}
	// crate transaction
	var tx *sql.Tx
	tx, err = edb.db.Begin()
	if err != nil {
		return err
	}
	if en.EntryID == 0 {
		// insert records
		err = en.dbInsert(tx)
		if err != nil {
			goto DONE
		}

		es.EntryFK = en.EntryID
		err = es.dbInsert(tx)
		if err != nil {
			goto DONE
		}

		if r != nil {
			r.EntryFK = en.EntryID
			err = r.dbInsert(tx)
			// if err != nil {
			// 	goto DONE
			// }
		}
	} else {
		// update records
		err = en.dbUpdate(tx)
		if err != nil {
			goto DONE
		}
		err = es.dbUpdate(tx)
	}
DONE:
	// check transaction status
	if err != nil {
		// rollback due to error
		log.Printf("Failed to save. Error: %v", err)
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("Failed to rollback. Error: %v", err2)
		}
		return err
	}
	// commit
	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("Failed to commit entry to DB. Error: %v", err)
	}
	return err
}

// RecentHTMLEntries returns limit recent entries with HTML content
// ordered in DESC order.  The data is suitable for viewing, but not editing.
func (edb *EntryDB) RecentHTMLEntries(limit int64) (result []WSEntryGetHTML, err error) {
	if edb.db == nil {
		return result, fmt.Errorf("Database connection is closed.")
	}
	var rows *sql.Rows
	sql := `
	SELECT e.entryID, es.title, e.titleIcon, e.html, e.rawType, e.updated
	from entry e
	inner join entrySearch es on es.EntryFK = e.EntryID
	order by e.updated desc limit ?
	`
	rows, err = edb.db.Query(sql, limit)
	return edb.rowsToWSEntryGetHTML(rows, err)
}

// MatchEntries searches enrySearch table for matching entires.
func (edb *EntryDB) MatchEntries(query string, limit int64) (result []WSEntryGetHTML, err error) {
	if edb.db == nil {
		return result, fmt.Errorf("Database connection is closed.")
	}
	var rows *sql.Rows
	sql := `
SELECT e.entryID, es.title, e.titleIcon, e.html, e.rawType, e.updated from entry e
inner join entrySearch es on es.EntryFK = e.EntryID
where entryID in (select entryFK from entrySearch where entrySearch match $1)
order by updated desc limit $2
`
	rows, err = edb.db.Query(sql, query, limit)
	return edb.rowsToWSEntryGetHTML(rows, err)
}

// rowsToWSEntryGetHTML is a common results processing function
// to return list of entries to view.
func (edb *EntryDB) rowsToWSEntryGetHTML(rows *sql.Rows, err error) ([]WSEntryGetHTML, error) {
	if err != nil {
		return nil, err
	}
	result := []WSEntryGetHTML{}
	var r WSEntryGetHTML
	for rows.Next() {
		var rawType int
		err = rows.Scan(&r.EntryID, &r.Title, &r.TitleIcon, &r.HTML, &rawType, &r.Updated)
		if err != nil {
			return result, err
		}
		// resolve number to a name
		r.RawTypeName = edb.rawTypes.NameByNum(rawType)
		result = append(result, r)
	}
	return result, err
}

// GetFullEntry loads all Entry data for editing.
func (edb *EntryDB) GetFullEntry(entryID int64) (r *WSFullEntry, err error) {
	r = &WSFullEntry{EntryID: entryID}
	if edb.db == nil {
		return r, fmt.Errorf("Database connection is closed.")
	}
	sql := `
	SELECT es.title, e.titleIcon, e.rawType, e.raw, e.html, e.updated, es.tags
	from entry e
	inner join entrySearch es on e.entryID = es.entryFK
	where e.entryID = ?
	`
	var rawType int
	err = edb.db.QueryRow(sql, entryID).
		Scan(&r.Title, &r.TitleIcon, &rawType, &r.Raw, &r.HTML, &r.Updated, &r.Tags)
	if err != nil {
		return r, err
	}
	r.RawTypeName = edb.rawTypes.NameByNum(rawType)
	if r.RawTypeName == "" {
		return r, fmt.Errorf("Error loading RawTypeName for number %v", rawType)
	}
	return r, nil
}

// GetOrCreateUser gets existing user or creates new one with basic (Guest) clearance.
func (edb *EntryDB) GetOrCreateUser(gUser goth.User) (user *User, err error) {
	if edb.db == nil {
		return user, fmt.Errorf("Database connection is closed.")
	}
	query := `
	SELECT u.userID, clearances from user u
	inner join oauthUser ou on ou.UserFK = u.UserID
	where ou.provider = ? and ou.provUserID = ?
	`
	user = &User{}
	err = edb.db.QueryRow(query, gUser.Provider, gUser.UserID).
		Scan(&user.UserID, &user.Clearances)
	if err == sql.ErrNoRows {
		return edb.createUser(gUser)
	}
	return
}

func (edb *EntryDB) createUser(gUser goth.User) (user *User, err error) {
	var tx *sql.Tx
	tx, err = edb.db.Begin()
	if err != nil {
		return nil, err
	}
	user = &User{}
	ou := OAuthUser{
		Provider:          gUser.Provider,
		Email:             gUser.Email,
		Name:              gUser.Name,
		FirstName:         gUser.FirstName,
		LastName:          gUser.LastName,
		NickName:          gUser.NickName,
		Description:       gUser.Description,
		ProvUserID:        gUser.UserID,
		AvatarURL:         gUser.AvatarURL,
		Location:          gUser.Location,
		AccessToken:       gUser.AccessToken,
		AccessTokenSecret: gUser.AccessTokenSecret,
		RefreshToken:      gUser.RefreshToken,
		ExpiresAt:         gUser.ExpiresAt,
	}
	err = user.dbInsert(tx)
	if err != nil {
		goto DONE
	}
	ou.UserFK = user.UserID
	err = ou.dbInsert(tx)
DONE:
	// check transaction status
	if err != nil {
		// rollback due to error
		log.Printf("Failed to save. Error: %v", err)
		err2 := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback. Error: %v", err2)
		}
		return nil, err
	}
	// commit
	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("Failed to commit entry to DB. Error: %v", err)
	}
	return user, err
}

// GetUser loads User by userID.
func (edb *EntryDB) GetUser(userID int64) (u *WSUserGet, err error) {
	u = nil
	if edb.db == nil {
		err = errors.New("Database connection is closed.")
		return
	}
	query := `SELECT u.clearances, ou.name, ou.nickName, ou.avatarURL from user u
	inner join OAuthUser ou on ou.UserFK = u.UserID
	where u.userId = ?`
	u = &WSUserGet{}
	err = edb.db.QueryRow(query, userID).Scan(&u.Clearances, &u.Name, &u.NickName, &u.AvatarURL)
	return
}

// GetUsers returns web safe user list.
func (edb *EntryDB) GetUsers() (users []*WSUserGet, err error) {
	if edb.db == nil {
		err = errors.New("Database connection is closed.")
		return
	}
	query := `
	SELECT u.clearances, ou.name, ou.nickName, ou.avatarURL from user u
	inner join OAuthUser ou on ou.UserFK = u.UserID
	order by u.Updated desc
	`
	var rows *sql.Rows
	rows, err = edb.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := WSUserGet{}
		err = rows.Scan(&user.Clearances, &user.Name, &user.NickName, &user.AvatarURL)
		if err != nil {
			return
		}
		users = append(users, &user)
	}
	return
}

// GetUniqueRedirectPaths returns all unique 1st hop paths.
// For example, /hop1/hop2/hop2 results in just /hop1
// If original path starts with slash, then result contains /
// If original path does not start with slash, then result is void of it.
func (edb *EntryDB) GetUniqueRedirectPaths() (result []string, err error) {
	if edb.db == nil {
		return result, fmt.Errorf("Database connection is closed.")
	}
	var rows *sql.Rows
	sql := `SELECT path from redirect`
	rows, err = edb.db.Query(sql)
	if err != nil {
		return nil, err
	}
	result = []string{}
	unique := map[string]string{}
	var path string
	for rows.Next() {
		err = rows.Scan(&path)
		if err != nil {
			return
		}
		idx := 0
		if strings.HasPrefix(path, "/") {
			idx = 1
		}
		// string "a/b" is split into "a", "b"
		// string "/a/b" is split into "", "a", "b"
		// thus if it starts with /, then 1st is ""
		// we are looking for 1st meaningful prefix and that's a
		hops := strings.Split(path, "/")
		prefix := hops[idx]
		if idx == 1 {
			// add slash back, so it is in unique list
			prefix = "/" + prefix
		}
		_, ok := unique[prefix]
		if !ok {
			unique[prefix] = ""
			result = append(result, prefix)
		}
	}
	return
}

// GetRedirectEntryID returns EntryID corresponding to redirect item.
func (edb *EntryDB) GetRedirectEntryID(path string) (entryID int64, err error) {
	if edb.db == nil {
		err = fmt.Errorf("Database connection is closed.")
		return
	}
	query := "select entryFK from redirect where Path=?"
	err = edb.db.QueryRow(query, path).Scan(&entryID)
	return
}
