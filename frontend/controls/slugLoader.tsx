/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import React from 'react'
import {WSEntryGetHTML} from '../model/wsentry'
import {DataService} from '../common/dataService'
import {EntryViewBox} from '../controls/entryView'
import * as router from 'react-router'

declare var ZonePath: string

export interface SlugLoaderProps {
	params: router.Params
}

class SlugLoaderState {
	constructor(
		public entry: WSEntryGetHTML,
	) {}
}

export class SlugLoaderBox extends React.Component<SlugLoaderProps, SlugLoaderState> {
	public constructor(props: SlugLoaderProps) {
		super(props)
		let slug = props.params['slug'] as string
		this.state = new SlugLoaderState(null);
		DataService.get(ZonePath + '/api/s/' + slug)
		.then(function(jsonEntry) {
			console.log("loaded jsonEntry by slug", jsonEntry)
			this.setState(new SlugLoaderState(jsonEntry as WSEntryGetHTML))
		}.bind(this))
		.catch(function(err) {
			console.log("err loading json: ", err)
		}.bind(this))
	}

	render() {
		//console.log("SlugLoader state", this.state)
		if (this.state.entry == null) {
			return (<p>Loading entry...</p>)
		} else {
			return (<EntryViewBox entry={this.state.entry} expandInitially={true} />)
		}
	} // end of render function
} // end of class
