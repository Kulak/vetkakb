/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import * as React from 'react'
import {WSEntryGetHTML} from '../model/wsentry'
import {DataService} from '../common/dataService'
import {EntryViewBox} from '../controls/entryView'
import * as router from 'react-router'

declare var ZonePath: string

export interface EntryLoaderProps {
	params: router.Params
}

class EntryLoaderState {
	constructor(
		public entry: WSEntryGetHTML,
	) {}
}

export class EntryLoaderBox extends React.Component<EntryLoaderProps, EntryLoaderState> {
	public constructor(props: EntryLoaderProps) {
		super(props)
		let entryID = props.params['entryID'] as string
		this.state = new EntryLoaderState(null);
		DataService.get(ZonePath + '/api/entry/' + entryID)
		.then(function(jsonEntry) {
			console.log("loaded jsonEntry", jsonEntry)
			this.setState(new EntryLoaderState(jsonEntry as WSEntryGetHTML))
		}.bind(this))
		.catch(function(err) {
			console.log("err loading json: ", err)
		}.bind(this))
	}

	// onEditorCloseRequested(fe: WSFullEntry) {
	// 	console.log("entryLoader: set new EntryViewState", fe)
	// 	this.setState(new EntryLoaderState(fe, true, false))
	// }

	render() {
		//console.log("entryLoader state", this.state)
		if (this.state.entry == null) {
			return (<p>Loading entry...</p>)
		} else {
			return (<EntryViewBox entry={this.state.entry} />)
		}
	} // end of render function
} // end of class
