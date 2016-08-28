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
