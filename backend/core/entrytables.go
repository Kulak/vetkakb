package core

import (
	"database/sql"
	"fmt"
	"time"
)

// Entry represents content of Entry table in Entry databse.
type Entry struct {
	EntryID           int64
	Raw               []byte
	RawType           int
	RawContentType    string
	RawFileName       string
	HTML              string
	OwnerFK           int64
	RequiredClearance uint8
	Created           time.Time
	Updated           time.Time
}

// EntrySearch represents content of EntrySearch in Entry databse.
type EntrySearch struct {
	// EntryFK is a foreign key into EntryID of Entry.
	EntryFK int64
	// Title is entry title.
	Title string
	// Plain represents indexed content of the entry.
	Plain string
	// Tags is a comma separated list of tags.
	Tags string
}

/* ========== Shared Functions ========== */

func sqlRequireAffected(result sql.Result, expected int64) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("Expect affected DB rows to be %v, but was %v", expected, affected)
	}
	return nil
}

/* ========== Entry ========== */

// NewEntry creates new entry to be inserted into DB.
func NewEntry(entryID int64, raw []byte, rawType int, rawContentType,
	rawFileName string, ownerFK int64) *Entry {
	return &Entry{
		EntryID:        entryID,
		Raw:            raw,
		RawType:        rawType,
		RawContentType: rawContentType,
		RawFileName:    rawFileName,
		OwnerFK:        ownerFK,
	}
}

// savdbInserteToDB inserts record into DB.  EntryID must be zero.
// If operation is successful EntryID is set to inserted record.
func (en *Entry) dbInsert(tx *sql.Tx) (err error) {
	var result sql.Result
	var sql string
	if en.EntryID != 0 {
		return fmt.Errorf("Cannot insert record with existing EntryID %v.", en.EntryID)
	}
	sql = `
insert into entry (raw, rawType, rawContentType, rawFileName, html, ownerFK)
values($1, $2, $3, $4, $5, $6)
`
	result, err = tx.Exec(sql,
		en.Raw, en.RawType, en.RawContentType, en.RawFileName, en.HTML, en.OwnerFK)
	if err != nil {
		return fmt.Errorf("Failed to insert `entry` record to DB. Error: %v", err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf(
			"Failed to get affected `entry` records after insert to DB. Error: %v",
			err)
	}
	en.EntryID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to get EntryID for last insert operation. Error: %v", err)
	}
	return err
}

// dbUpdate updates existing record.
func (en *Entry) dbUpdate(tx *sql.Tx) (err error) {
	var result sql.Result
	var sql string
	if en.EntryID == 0 {
		return fmt.Errorf("Cannot update record, because EntryID is set to zero.")
	}
	sql = "update `entry` set raw=$1, rawType=$2, rawContentType=$3, rawFileName=$4, html=$5, updated=$6 " +
		"where entryID=$7"
	en.Updated = time.Now()
	result, err = tx.Exec(sql,
		en.Raw, en.RawType, en.RawContentType, en.RawFileName, en.HTML, en.Updated.Unix(),
		en.EntryID)
	if err != nil {
		return fmt.Errorf("Failed to update EntryID %v. Error: %v", en.EntryID, err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf(
			"Failed to get affected records after EntryID %v update. Error: %v",
			en.EntryID, err)
	}
	return err
}

/* ========== EntrySearch ========== */

// NewEntrySearch creates new entry search item to be inserted into DB.
func NewEntrySearch(entryFK int64, title string, tags string) *EntrySearch {
	return &EntrySearch{
		EntryFK: entryFK,
		Title:   title,
		Tags:    tags,
	}
}

func (es *EntrySearch) dbInsert(tx *sql.Tx) (err error) {
	var result sql.Result
	var sql string
	sql = "insert into `entrySearch` (entryFK, title, plain, tags) values($1, $2, $3, $4)"
	result, err = tx.Exec(sql, es.EntryFK, es.Title, es.Plain, es.Tags)
	if err != nil {
		return fmt.Errorf("Failed to insert `entrySeach` record to DB. Error: %v", err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf(
			"Failed to get affected `entrySearch` records after insert. Error: %v",
			err)
	}
	var insertedID int64
	insertedID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf(
			"Failed to ID of last insert operation on entrySearch table. Error: %v",
			err)
	}
	if insertedID != es.EntryFK {
		return fmt.Errorf("Inserted entrySearch ID %v does not match expected %v",
			insertedID, es.EntryFK)
	}
	return err
}

func (es *EntrySearch) dbUpdate(tx *sql.Tx) (err error) {
	var result sql.Result
	var sql string
	sql = "update `entrySearch` set title=$1, plain=$2, tags=$3 where EntryFK=$4"
	result, err = tx.Exec(sql, es.Title, es.Plain, es.Tags, es.EntryFK)
	if err != nil {
		return fmt.Errorf("Failed to update `entrySeach` record. Error: %v", err)
	}
	err = sqlRequireAffected(result, 1)
	if err != nil {
		return fmt.Errorf(
			"Failed after update check of `entrySearch` affected records for EntryFK %v. Error: %v",
			es.EntryFK, err)
	}
	return err
}
