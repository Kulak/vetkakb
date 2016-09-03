/*
Runs search query and displays results in the list.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryList} from './entryList'

export class SearchProps {}

class SearchState {
	constructor(
		public entries: Array<WSEntryGetHTML> = []
	) {}
}

export class Search extends React.Component<SearchProps, SearchState> {

	public constructor(props: SearchProps, context: any) {
		super(props, context)
		this.state = new SearchState()
		DataService.get('/api/search/')
		.then(function(jsonEntries) {
			console.log("search results", jsonEntries)
			this.setState(new SearchState(jsonEntries as Array<WSEntryGetHTML>))
		}.bind(this))
		.catch(function(err) {
			console.log("search error: ", err)
		}.bind(this))
	}

	render() {
		return <EntryList entries={this.state.entries} />
	}
}  // end of class