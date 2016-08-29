/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import * as React from 'react'
import {WSEntryGetHTML} from '../model/wsentry'

export class EntryViewProps {
	constructor(
		public entry: WSEntryGetHTML
	){}
}

class EntryViewState {
	constructor(
		public expanded: boolean = false,
		public editing: boolean = false
	) {}
}

export class EntryViewBox extends React.Component<EntryViewProps, EntryViewState> {
	public constructor(props: EntryViewProps, context) {
		super(props, context)
		this.state = new EntryViewState();
	}
	onExpandClick(expandAction: boolean) {
		this.setState(new EntryViewState(expandAction, false))
	}
	onEditClick(editAction: boolean) {
		this.setState(new EntryViewState(false, editAction))
	}
	render() {
		let en = this.props.entry
		if (this.state.editing) {

		} else {
			if (this.state.expanded) {
				return <div>
					<h2 onClick={e => this.onExpandClick(false)}>{en.Title}
						<button onClick={e => this.onEditClick(true)}>Change</button>
					</h2>
					<div dangerouslySetInnerHTML={{__html: en.HTML}} />
				</div>
			} else {
				return <div><h2 onClick={e => this.onExpandClick(true)}>{en.Title}</h2></div>
			}
		}
	} // end of render function
} // end of class