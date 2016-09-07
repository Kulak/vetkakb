
const RTTextPlain: string = "Plain Text"

export class Entry {
	public constructor(
		public entryID: number = 0,
		public title: string = "",
		public raw: string = "",
		public rawTypeName: string = RTTextPlain,
		public tags: string = ""
	) {}
}