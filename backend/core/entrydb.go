package core

import (
	"database/sql"
	"fmt"
	"log"
)

// EntryDB is a service to interact with EntryDB database.
type EntryDB struct {
	db         *sql.DB
	dbFileName string
}

// NewEntryDB creates DB service, but does not open connection.
func NewEntryDB(dbFileName string) *EntryDB {
	return &EntryDB{
		dbFileName: dbFileName,
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
		err = tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback. Error: %v", err)
		}
	} else {
		// commit
		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("Failed to commit entry to DB. Error: %v", err)
		}
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
	sql := "SELECT entryID, title, html, updated from `entry` order by updated desc limit ?"
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
SELECT entryID, title, html, updated from entry where entryID in
(select entryFK from entrySearch where entrySearch match $1)
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
		err = rows.Scan(&r.EntryID, &r.Title, &r.HTML, &r.Updated)
		if err != nil {
			return result, err
		}
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
	SELECT e.title, e.rawType, e.raw, e.html, e.updated, es.tags
	from entry e
	inner join entrySearch es on e.entryID = es.entryFK
	where e.entryID = ?
	`
	err = edb.db.QueryRow(sql, entryID).
		Scan(&r.Title, &r.RawType, &r.Raw, &r.HTML, &r.Updated, &r.Tags)
	return r, err
}
