/*
Displays an array of entires passed as parameters.
*/

import * as React from 'react'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryViewBox} from './entryView'

export interface EntryListProps {
	entries: Array<WSEntryGetHTML>
}

class EntryListState {
}

export class EntryList extends React.Component<EntryListProps, EntryListState> {

	public constructor(props: EntryListProps, context: any) {
		super(props, context)
		this.state = new EntryListState()
	}

	render() {
		let entries = this.props.entries.map(function(entry: WSEntryGetHTML) {
			return <EntryViewBox entry={entry} expandInitially={false} key={entry.EntryID}  />
		})
		return <div>
			{entries}
		</div>
	}
}  // end of class
