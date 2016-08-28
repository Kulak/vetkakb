
const RTTextPlain: number = 1

export class Entry {
	public constructor(
		public entryID: number = 0,
		public title: string = "",
		public raw: string = "",
		public rawType: number = RTTextPlain,
		public tags: string = ""
	) {}
}