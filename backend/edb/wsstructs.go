package edb

import (
	"fmt"
	"time"
)

/*
This file collects data structures that support REST API calls.
*/

// WSUserGet is a structure required by HTML app
// to display basic user data.
type WSUserGet struct {
	Name       string
	Clearances uint8
	NickName   string
	AvatarURL  string
}

// GuestWSUserGet represents default Guest user.
var GuestWSUserGet = &WSUserGet{
	Name:       "Anonymous Guest",
	Clearances: Guest.Mask,
	NickName:   "Guest",
	AvatarURL:  "",
}

// WSEntryPut describes REST API payload for creating entry.
type WSEntryPut struct {
	Title          string
	TitleIcon      string
	RawTypeName    string
	RawContentType string
	RawFileName    string
	Tags           string
	Intro          string
	Slug           string
}

// WSEntryPost is used to udpate enries.
type WSEntryPost struct {
	EntryID int64
	WSEntryPut
}

// WSEntryGetHTML is used to load a list of entries to display text and HTML oriented content.
// This structure is populated directly by DAL.
type WSEntryGetHTML struct {
	EntryID     int64
	Title       string
	TitleIcon   string
	RawTypeName string
	HTML        string
	Intro       string
	Slug        string
	Updated     time.Time
}

// WSFullEntry is used to load detailed Entry data used in entry editor.
type WSFullEntry struct {
	EntryID     int64
	Title       string
	TitleIcon   string
	RawTypeName string
	Raw         []byte
	HTML        string
	Intro       string
	Tags        string
	Slug        string
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
