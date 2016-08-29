

// WSEntryPut mirrors backend structure with the same name.
export class WSEntryPut {
	constructor(
		public title: string,
		// raw shall be base64 encoded
		public raw: string,
		public rawType: number,
		public tags: string) {}
}