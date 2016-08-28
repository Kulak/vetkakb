package vetka

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

// WSEntryGet is used to load entries.
type WSEntryGet struct {
}
