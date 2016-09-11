package core

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/markbates/goth"
)

// EntryDB is a service to interact with EntryDB database.
type EntryDB struct {
	db         *sql.DB
	dbFileName string
	rawTypes   *TypeService
}

// NewEntryDB creates DB service, but does not open connection.
func NewEntryDB(dbFileName string, rawTypes *TypeService) *EntryDB {
	return &EntryDB{
		dbFileName: dbFileName,
		rawTypes:   rawTypes,
	}
}

// Open opens DB connection.
func (edb *EntryDB) Open() error {
	if edb.db != nil {
		return nil
	}
	db, err := sql.Open("sqlite3", edb.dbFileName)
	if err != nil {
		return fmt.Errorf("Failed to open database %v. Error: %v", edb.dbFileName, err)
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
func (edb *EntryDB) SaveEntry(en *Entry, es *EntrySearch) (err error) {
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
		if err != nil {
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
	SELECT e.entryID, es.title, e.html, e.rawType, e.updated from entry e
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
SELECT e.entryID, es.title, e.html, e.rawType, e.updated from entry e
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
		err = rows.Scan(&r.EntryID, &r.Title, &r.HTML, &rawType, &r.Updated)
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
	SELECT es.title, e.rawType, e.raw, e.html, e.updated, es.tags
	from entry e
	inner join entrySearch es on e.entryID = es.entryFK
	where e.entryID = ?
	`
	var rawType int
	err = edb.db.QueryRow(sql, entryID).
		Scan(&r.Title, &rawType, &r.Raw, &r.HTML, &r.Updated, &r.Tags)
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
	if edb.db == nil {
		return u, fmt.Errorf("Database connection is closed.")
	}
	query := `SELECT u.clearances, ou.name, ou.nickName, ou.avatarURL from user u
	inner join OAuthUser ou on ou.UserFK = u.UserID
	where u.userId = ?`
	u = &WSUserGet{}
	err = edb.db.QueryRow(query, userID).Scan(&u.Clearances, &u.Name, &u.NickName, &u.AvatarURL)
	return
}
