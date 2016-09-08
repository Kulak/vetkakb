import {DataService} from './dataService'

export class WSRawType {

	public static Binary = "Binary/"
	public static BinaryImage = "Binary/Image"
	public static Markdown = "Markdown"

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
						WSRawType._list.sort((a,b) => { return a.TypeNum - b.TypeNum })
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

	public static NameForNum(typeNum: number, list: Array<WSRawType>): string {
		let found = list.find(each => {
			return each.TypeNum == typeNum
		})
		if (found != null) {
			return found.Name
		}
		return ""
	}

}