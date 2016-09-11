
// WSEntryPut mirrors backend structure with the same name.
export class WSUserGet {
	constructor(
		public Name: string,
		public Clearances: number,
		public NickName: string,
		public AvatarURL: string) {}
}