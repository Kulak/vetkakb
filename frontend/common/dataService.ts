
// DataService provides functions to simplify
// and centralize interaction with REST API.
export class DataService {

	// tnewGetRequest is targeting GET requests 1st
	public static newGetRequest(url: string): Request {
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

	// newRequestWith is targeting PUT and POST requests 1st
	public static newRequestInit(method: string, data: any): RequestInit {
		// About fetch API and Request object as part of fetch:
		// https://developer.mozilla.org/en-US/docs/Web/API/Request/Request
		let payload: string = JSON.stringify(data)
		return {
			method: method,
			// with 'include' basic authentication header is sent with request
			//credentials: 'include'
			credentials: 'same-origin',
			headers: {
				'Accept': 'application/json',
				'Content-type': 'application/json'
			},
			body: payload
		}
	}
	// handleFetch executes fetch and treats 404 response code as an error.
	public static handleFetch<T>(url: string, reqInit: RequestInit): Promise<T> {
		return fetch(url, reqInit)
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

	// get function executes GET method and treats 404 response code as an error.
	public static get<T>(url: string): Promise<T> {
		return DataService.handleFetch(url, DataService.newGetRequest(url))
	}

	// put creates a PUT request and treates 404 as an error.
	public static put<T>(url: string, data: any): Promise<T> {
		let req = DataService.newRequestInit("PUT", data)
		return DataService.handleFetch(url, req)
	} // end of put

	// post creates a PUT request and treates 404 as an error.
	public static post<T>(url: string, data: any): Promise<T> {
		let req = DataService.newRequestInit("POST", data)
		return DataService.handleFetch(url, req)
	} // end of put

}  // end of class