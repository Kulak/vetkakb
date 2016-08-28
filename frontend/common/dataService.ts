
// DataService provides functions to simplify
// and centralize interaction with REST API.
export class DataService {

	// that's new GET request
	public static newRequest(url: string): Request {
		// About fetch API and Request object as part of fetch:
		// https://developer.mozilla.org/en-US/docs/Web/API/Request/Request
		return new Request(
			url, {
        // with 'include' basic authentication header is sent with request
        //credentials: 'include'
        credentials: 'same-origin'
			}
		)
	}

	// get function executes GET method and treats 404 response code as an error.
	public static get<T>(url: string): Promise<T> {
		return DataService.getr(DataService.newRequest(url))
	}

	// get function executes GET method and treats 404 response code as an error.
	public static getr<T>(req: Request): Promise<T> {
		return fetch(req)
  	.then(function(response) {
			// 404 code is a good response, so we need to check
			if (response.ok) {
				// continue chain
    		return response.json()
			} else {
				return response.text()
				.then(function(text) {
					throw new Error(text)
				})
			}
  	})
	} // end of getr<T>

	public static put<T>(url: string, data: any): Promise<T> {
		let payload: string = JSON.stringify(data)
		return fetch(url, {
			method: 'PUT',
			credentials: 'same-origin',
			headers: {
				'Accept': 'application/json',
				'Content-type': 'application/json'
			},
			body: payload
		})
  	.then(function(response) {
			// 404 code is a good response, so we need to check
			if (response.ok) {
				// continue chain
    		return response.json()
			} else {
				return response.text()
				.then(function(text) {
					throw new Error(text)
				})
			}
  	})
	}

}  // end of class