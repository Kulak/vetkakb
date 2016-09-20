package edb

import (
	"database/sql"
	"fmt"
)

// Redirect represents content of Redirect table.
type Redirect struct {
	RedirectID int64
	Path       string
	EntryFK    int64
	StatusCode int
}

func (r *Redirect) dbInsert(tx *sql.Tx) (err error) {
	var result sql.Result
	var query string
	if r.RedirectID != 0 {
		return fmt.Errorf("Cannot insert record with non-zero RedirectID %v", r.RedirectID)
	}
	query = `
insert into redirect (Path, EntryFK, StatusCode)
values($1, $2, $3)
`
	result, err = tx.Exec(query, r.Path, r.EntryFK, r.StatusCode)
	if err != nil {
		return fmt.Errorf("Failed to insert into redirect table.  Error: %s", err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf("Failed to insert into redirect table. Error: %v", err)
	}
	r.RedirectID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to get EntryID for last insert operation. Error: %v", err)
	}
	return nil
}
