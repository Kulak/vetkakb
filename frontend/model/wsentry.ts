

// RawTypes defines data types that can be stored.
export enum RawType {
	Undefined = 0,
	PlainText
}

// WSEntryPut mirrors backend structure with the same name.
export class WSEntryPut {
	constructor(
		public title: string,
		// raw shall be base64 encoded
		public raw: string,
		public rawType: number,
		public tags: string) {}
}

// WSEntryPut mirrors backend structure with the same name.
export class WSEntryPost {
	constructor(
		public entryID: number,
		public title: string,
		// raw shall be base64 encoded
		public raw: string,
		public rawType: number,
		public tags: string) {}
}

// WSEntryGetHTML mirrows backend structure
// and is used to load view only form of entry.
export class WSEntryGetHTML {
	constructor(
		public EntryID: number = 0,
		public Title: string = "",
		public HTML: string = "",
		public Updated: string = ""
	) {}
}

// WSFullEntry is returned by GET entry call.
// It is used by entryEditor.
export class WSFullEntry {
	constructor(
		public EntryID: number = 0,
		public Title: string = "",
		public Raw: string = "",
		public RawType: number = 0,
		public Tags: string = "",
		public HTML: string = "",
		public Updated: string = ""
	) {}
}