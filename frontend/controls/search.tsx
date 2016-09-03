/*
Search is a top level control that combines search conditions
with action to load data based on selected conditions
and with display of search results.

It is a core search controller.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryViewBox} from './EntryView'

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
	}

	onRecentClick() {
		console.log("in onRecentClick")

		DataService.get('/api/recent/20')
		.then(function(jsonEntries) {
			console.log("json text", jsonEntries)
			this.setState(new SearchState(jsonEntries as Array<WSEntryGetHTML>))
		}.bind(this))
		.catch(function(err) {
			console.log("err loading json: ", err)
		}.bind(this))
	}

	render() {
		let entries = this.state.entries.map(function(entry: WSEntryGetHTML) {
			console.log("map entryId, title", entry.EntryID, entry.Title)
			return <EntryViewBox entry={entry} key={entry.EntryID}  />
		})
		return <div>
			<button onClick={e => this.onRecentClick()}>Recent</button>
			{entries}
		</div>
	}
}  // end of class