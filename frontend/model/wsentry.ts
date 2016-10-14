
declare var ZonePath: string

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
		public Intro: string,
		public Slug: string,
	) {}
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
		public Slug: string = "",
		public Updated: string = "",
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
		public Slug: string = "",
		public Updated: string = "",
	) {}

	initializeFromWSFullEntry(src: WSFullEntry): WSFullEntry {
		this.EntryID = src.EntryID
		this.HTML = src.HTML
		this.Intro = src.Intro
		this.Raw = src.Raw
		this.RawTypeName = src.RawTypeName
		this.Tags = src.Tags
		this.Title = src.Title
		this.TitleIcon = src.TitleIcon
		this.Slug = src.Slug
		this.Updated = src.Updated
		return this
	}

	initializeFromWSEntryGetHTML(src: WSEntryGetHTML): WSFullEntry {
		this.EntryID = src.EntryID
		this.HTML = src.HTML
		this.Intro = src.Intro
		this.Raw = null
		this.RawTypeName = src.RawTypeName
		this.Tags = ""
		this.Title = src.Title
		this.TitleIcon = src.TitleIcon
		this.Slug = src.Slug
		this.Updated = src.Updated
		return this
	}

	copy(): WSFullEntry {
		return new WSFullEntry().initializeFromWSFullEntry(this)
	}

	// fromData is meant to load from REST API response.
	// The incoming data type is not trully a WSFullEntry,
	// because it has no functions associated with it.
	static fromData(data:WSFullEntry): Promise<WSFullEntry> {
		return new Promise((fulfil, reject) => {
			let r = new WSFullEntry().initializeFromWSFullEntry(data)
			// don't convert null, because it atob(null) returns "ée"
			if (r.Raw != null) {
				let blob = this.b64toBlob(r.Raw)
				let reader = new FileReader()
				reader.addEventListener('loadend', function() {
					// listener is called when readAsText is completed; promise could be used here
					// For ISO-8859-1 there's no further conversion required
					r.Raw = reader.result
					fulfil(r)
				})
				reader.readAsText(blob)
			} else {
				r.Raw = ""
				fulfil(r)
			}
		})
	}

	static b64toBlob (b64Data, contentType='', sliceSize=512) {
		const byteCharacters = atob(b64Data);
		const byteArrays = [];

		for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
			const slice = byteCharacters.slice(offset, offset + sliceSize);

			const byteNumbers = new Array(slice.length);
			for (let i = 0; i < slice.length; i++) {
				byteNumbers[i] = slice.charCodeAt(i);
			}

			const byteArray = new Uint8Array(byteNumbers);

			byteArrays.push(byteArray);
		}

		const blob = new Blob(byteArrays, {type: contentType});
		return blob;
	}

	permalink(): string {
		return ZonePath+'/app/e/'+this.EntryID
	}
}