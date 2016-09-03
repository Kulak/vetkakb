/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import * as React from 'react'
import {WSEntryGetHTML, WSFullEntry} from '../model/wsentry'
import {EntryEditor, EditorProps} from './EntryEditor'
import {Entry} from '../model/entry'
import {DataService} from '../common/dataService'

export interface EntryViewProps {
		entry: WSEntryGetHTML
}

class EntryViewState {
	constructor(
		public fullEntry: WSFullEntry,
		public expanded: boolean,
		public editing: boolean,
	) {}
}

export class EntryViewBox extends React.Component<EntryViewProps, EntryViewState> {
	public constructor(props: EntryViewProps) {
		super(props)
		let pe = props.entry
		let fe = new WSFullEntry(pe.EntryID, pe.Title, null, 0, "", pe.HTML, pe.Updated)
		this.state = new EntryViewState(fe, false, false);
	}
	onExpandClick(expandAction: boolean) {
		this.setState(new EntryViewState(this.state.fullEntry, expandAction, false))
	}
	onEditClick(editAction: boolean) {
		if (editAction) {
			// load a full entry
			DataService.get('/api/entry/' + this.props.entry.EntryID)
			.then(function(jsonEntry) {
				console.log("json text", jsonEntry)
				let entry = jsonEntry as WSFullEntry
				entry.Raw = atob(entry.Raw)
				this.setState(new EntryViewState(jsonEntry as WSFullEntry, false, editAction))
			}.bind(this))
			.catch(function(err) {
				console.log("err loading json: ", err)
			}.bind(this))
		} else {
			this.setState(new EntryViewState(this.state.fullEntry, false, editAction))
		}
	}
	onEditorCloseRequested(fe: WSFullEntry) {
		console.log("entryView: set new EntryViewState", fe)
		this.setState(new EntryViewState(fe, true, false))
	}
	render() {
		let fe: WSFullEntry = this.state.fullEntry
		console.log("entryView: render entry", fe)
		if (this.state.editing) {
			return <EntryEditor entry={fe} editorCloseReq={fe => this.onEditorCloseRequested(fe)} />
		} else {
			if (this.state.expanded) {
				return <div>
					<div className='toolbar entryHeader'>
						<h2 className='leftStack' onClick={e => this.onExpandClick(false)}>{fe.Title}</h2>
						<button className='leftStack' onClick={e => this.onEditClick(true)}>Change</button>
					</div>
					<div className='entryBody' dangerouslySetInnerHTML={{__html: fe.HTML}} />
				</div>
			} else {
				return <div><h2 className='entryHeader' onClick={e => this.onExpandClick(true)}>{fe.Title}</h2></div>
			}
		}
	} // end of render function
} // end of class