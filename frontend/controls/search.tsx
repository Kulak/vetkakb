/*
Runs search query and displays results in the list.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryList} from './entryList'

declare var ZonePath: string

export class SearchProps {}

class SearchState {
	constructor(
		public query: string = "",
		public entries: Array<WSEntryGetHTML> = []
	) {}
}

export class Search extends React.Component<SearchProps, SearchState> {

	public constructor(props: SearchProps, context: any) {
		super(props, context)
		this.state = new SearchState()
	}

	onQueryChange(event: React.FormEvent) {
		let query = (event.target as any).value
		this.setState(new SearchState(query))
	}

	onQuerySubmit(event: React.FormEvent) {
		event.preventDefault()
		this.sendQuery()
	}

	sendQuery() {
		let query = this.state.query
		DataService.get(ZonePath + '/api/search/' + query)
		.then(function(jsonEntries) {
			console.log("search results", jsonEntries)
			this.setState(new SearchState(query, jsonEntries as Array<WSEntryGetHTML>))
		}.bind(this))
		.catch(function(err) {
			console.log("search error: ", err)
		}.bind(this))
	}

	render() {
		// value={this.state.entry.Title}
		return <div>
			<div className='toolbar'>
				<label className='leftStack'>Search:</label>
				<form onSubmit={e => this.onQuerySubmit(e)} >
					<input className='leftStack' type="input"
						onChange={e => this.onQueryChange(e)} />
						<button className='leftStack' type="submit">Submit</button>
				</form>
			</div>
			<EntryList entries={this.state.entries} />
		</div>
	}
}  // end of class