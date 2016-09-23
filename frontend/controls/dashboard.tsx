/*
Loads most recent entries and displays in a list.
Based on Recent control.
*/

import * as React from 'react'
import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'
import {Recent} from './recent'
import {EntryList} from './entryList'

declare var ZonePath: string

export class DashboardProps {}

class DashboardState {
	constructor(
	) {}
}

export class Dashboard extends React.Component<DashboardProps, DashboardState> {

	public constructor(props: DashboardProps, context: any) {
		super(props, context)
		this.state = new DashboardState()
	}

	render() {
		return (<Recent limit={10} />)
	}
}  // end of class