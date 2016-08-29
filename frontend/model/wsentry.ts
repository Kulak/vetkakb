

// WSEntryPut mirrors backend structure with the same name.
export class WSEntryPut {
	constructor(
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
		public EntryID: number,
		public Title: string,
		public TawType: number,
		public HTML: string,
		public Updated: string
	) {}
}