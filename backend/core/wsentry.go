package core

import "time"

/*
This file collects data structures that support REST API calls.
*/

// WSEntryPut describes REST API payload for creating entry.
type WSEntryPut struct {
	Title   string
	Raw     []byte
	RawType int
	Tags    string
}

// WSEntryPost is used to udpate enries.
type WSEntryPost struct {
}

// WSEntryGetHTML is used to load entries to display text and HTML oriented content.
// This structure is populated directly by DAL.
type WSEntryGetHTML struct {
	EntryID int64
	Title   string
	HTML    string
	Updated time.Time
}

// WSFullEntry is used to load detailed Entry data used in entry editor.
type WSFullEntry struct {
	EntryID int64
	Title   string
	RawType int
	Raw     []byte
	HTML    string
	Tags    string
	Updated time.Time
}
