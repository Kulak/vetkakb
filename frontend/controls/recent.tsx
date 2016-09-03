/*
Loads most recent entries and displays in a list.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryList} from './entryList'

export class RecentProps {}

class RecentState {
	constructor(
		public entries: Array<WSEntryGetHTML> = []
	) {}
}

export class Recent extends React.Component<RecentProps, RecentState> {

	public constructor(props: RecentProps, context: any) {
		super(props, context)
		this.state = new RecentState()
		DataService.get('/api/recent/20')
		.then(function(jsonEntries) {
			console.log("json text", jsonEntries)
			this.setState(new RecentState(jsonEntries as Array<WSEntryGetHTML>))
		}.bind(this))
		.catch(function(err) {
			console.log("err loading json: ", err)
		}.bind(this))
	}

	render() {
		return <EntryList entries={this.state.entries} />
	}
}  // end of class