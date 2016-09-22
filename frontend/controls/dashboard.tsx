/*
Loads most recent entries and displays in a list.
Based on Recent control.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryList} from './entryList'

declare var ZonePath: string

export class DashboardProps {}

class DashboardState {
	constructor(
		public entries: Array<WSEntryGetHTML> = []
	) {}
}

export class Dashboard extends React.Component<DashboardProps, DashboardState> {

	public constructor(props: DashboardProps, context: any) {
		super(props, context)
		this.state = new DashboardState()
		DataService.get(ZonePath + '/api/recent/5')
		.then(function(jsonEntries) {
			console.log("json text", jsonEntries)
			this.setState(new DashboardState(jsonEntries as Array<WSEntryGetHTML>))
		}.bind(this))
		.catch(function(err) {
			console.log("err loading json: ", err)
		}.bind(this))
	}

	render() {
		return <EntryList entries={this.state.entries} />
	}
}  // end of class