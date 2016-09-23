

// RawTypes defines data types that can be stored.
export class RawType {
	public static Undefined = "Unknown"
	public static PlainText = "Plain Text"
}

// WSEntryPut mirrors backend structure with the same name.
// Post is used to both create new and update existing entries.
// New entries have entryID set to zero.
// Existing entries have entryID set to above zero.
export class WSEntryPost {
	constructor(
		public entryID: number,
		public title: string,
		public titleIcon: string,
		public rawTypeName: string,
		public tags: string,
		public Intro: string) {}
}

// WSEntryGetHTML mirrows backend structure
// and is used to load view only form of entry.
export class WSEntryGetHTML {
	constructor(
		public EntryID: number = 0,
		public Title: string = "",
		public TitleIcon: string = "",
		public HTML: string = "",
		public Intro: string = "",
		public RawTypeName: string = "",
		public Updated: string = ""
	) {}
}

// WSFullEntry is returned by GET entry call.
// It is used by entryEditor.
export class WSFullEntry {
	constructor(
		public EntryID: number = 0,
		public Title: string = "",
		public TitleIcon: string = "",
		public Raw: string = "",
		public RawTypeName: string = "",
		public Tags: string = "",
		public HTML: string = "",
		public Intro: string = "",
		public Updated: string = ""
	) {}
}