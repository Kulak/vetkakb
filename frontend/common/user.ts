import {DataService} from './dataService'
import {WSUserGet} from '../model/wsuser'

export class User {

	private static _current: WSUserGet;

	constructor(
	) {}

	// Returns a list of RawTypes.
	public static Current(): Promise<WSUserGet> {
		return new Promise(function(fulfil, reject) {
			if (User._current == null) {
				DataService.get("/api/session/user")
					.then(function(json) {
						User._current = json as WSUserGet
						fulfil(User._current)
					})
					.catch(function(err) {
						reject(err)
					})
			} else {
				fulfil(User._current)
			}
		})
	}
}