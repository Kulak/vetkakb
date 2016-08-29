/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import * as React from 'react'
//import {Entry} from '../model/entry'
//import {WSEntryPut} from '../model/wsentry'
//import {DataService} from '../common/dataService'
import {WSEntryGetHTML} from '../model/wsentry'

export class EntryViewProps {
	constructor(
		public entry: WSEntryGetHTML
	){}
}

class EntryViewState {
	constructor(
		public expanded: boolean = false
	) {}
}

export class EntryViewBox extends React.Component<EntryViewProps, EntryViewState> {
	public constructor(props: EntryViewProps, context) {
		super(props, context)
		this.state = new EntryViewState();
	}
	onExpandClick() {
		this.setState(new EntryViewState(true))
	}
	onContractClick() {
		this.setState(new EntryViewState(false))
	}
	render() {
		let en = this.props.entry
		if (this.state.expanded) {
			return <div>
				<h2 onClick={e => this.onContractClick()}>{en.Title}</h2>
				<div dangerouslySetInnerHTML={{__html: en.HTML}} />
			</div>
		} else {
			return <div><h2 onClick={e => this.onExpandClick()}>{en.Title}</h2></div>
		}
	} // end of render function
} // end of class