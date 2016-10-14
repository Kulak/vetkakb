/*
Loads most recent entries and displays in a list.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryList} from './entryList'

declare var ZonePath: string

export class RecentProps {
		limit: number = 20
}

class RecentState {
	constructor(
		public entries: Array<WSEntryGetHTML> = []
	) {}
}

export class Recent extends React.Component<RecentProps, RecentState> {

	public constructor(props: RecentProps, context: any) {
		super(props, context)
		this.state = new RecentState()
		let limit = 20
		if (props.limit) {
			limit = props.limit
		}
		this.loadRecent(limit, "")
	}

	loadRecent(limit: number, endStr: string) {
		DataService.get(ZonePath + '/api/recent/' + limit + "/" + endStr)
		.then(function(jsonEntries) {
			console.log("loaded jsonEntries in Recent", jsonEntries)
			this.setState(new RecentState(jsonEntries as Array<WSEntryGetHTML>))
		}.bind(this))
		.catch(function(err) {
			console.log("err loading json: ", err)
		}.bind(this))
	}

	onLoadClick() {
		console.log("state.entires", this.state)
		let idx = this.state.entries.length - 1
		let end = this.state.entries[idx]
		this.loadRecent(10, end.Updated)
	}

	render() {
		//console.log("recent entries in render", this.state.entries)
		return (<div>
				<EntryList entries={this.state.entries} />
				<button onClick={e => this.onLoadClick()}>More</button>
			</div>
		)
	}
}  // end of class