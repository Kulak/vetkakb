package core

import (
	"fmt"
	"time"
)

/*
This file collects data structures that support REST API calls.
*/

// WSEntryPut describes REST API payload for creating entry.
type WSEntryPut struct {
	Title          string
	RawTypeName    string
	RawContentType string
	RawFileName    string
	Tags           string
}

// WSEntryPost is used to udpate enries.
type WSEntryPost struct {
	EntryID int64
	WSEntryPut
}

// WSEntryGetHTML is used to load entries to display text and HTML oriented content.
// This structure is populated directly by DAL.
type WSEntryGetHTML struct {
	EntryID     int64
	Title       string
	RawTypeName string
	HTML        string
	Updated     time.Time
}

// WSFullEntry is used to load detailed Entry data used in entry editor.
type WSFullEntry struct {
	EntryID     int64
	Title       string
	RawTypeName string
	Raw         []byte
	HTML        string
	Tags        string
	Updated     time.Time
}

// **** Functions ****

func (eput WSEntryPut) String() string {
	return fmt.Sprintf("WSEntryPut {Title: %s, RawTypeName: %s, Tags: %s}",
		eput.Title, eput.RawTypeName, eput.Tags)
}

func (epost WSEntryPost) String() string {
	return fmt.Sprintf("WSEntryPost {EntryID: %d, %s}", epost.EntryID, epost.WSEntryPut.String())
}
