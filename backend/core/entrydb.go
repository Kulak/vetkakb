package core

import "time"

// Entry represents content of Entry table in Entry databse.
type Entry struct {
	EntryID int
	Title   string
	Raw     []byte
	RawTyep string
	Created time.Time
	Updated time.Time
}

// EntrySearch represents content of EntrySearch in Entry databse.
type EntrySearch struct {
	// EntryFK is a foreign key into EntryID of Entry.
	EntryFK int
	// Plain represents indexed content of the entry.
	Plain string
	// Tags is a comma separated list of tags.
	Tags string
}
