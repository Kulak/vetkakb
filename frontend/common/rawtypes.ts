import {DataService} from './dataService'

export class WSRawType {
	private static _list: Array<WSRawType>;

	constructor(
		public TypeNum: number,
		public Name: string
	) {}

	// Returns a list of RawTypes.
	public static List(): Promise<Array<WSRawType>> {
		return new Promise(function(fulfil, reject) {
			if (WSRawType._list == null) {
				DataService.get("/api/rawtype/list")
					.then(function(json) {
						WSRawType._list = json as Array<WSRawType>
						fulfil(WSRawType._list)
					})
					.catch(function(err) {
						reject(err)
					})
			} else {
				fulfil(WSRawType._list)
			}
		})
	}
}