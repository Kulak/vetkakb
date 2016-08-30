/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import * as React from 'react'
import {WSEntryGetHTML} from '../model/wsentry'
import {EntryEditor, EditorProps} from './EntryEditor'
import {Entry} from '../model/entry'

export interface EntryViewProps {
		entry: WSEntryGetHTML
}

class EntryViewState {
	constructor(
		public expanded: boolean = false,
		public editing: boolean = false
	) {}
}

export class EntryViewBox extends React.Component<EntryViewProps, EntryViewState> {
	public constructor(props: EntryViewProps) {
		super(props)
		this.state = new EntryViewState();
	}
	onExpandClick(expandAction: boolean) {
		this.setState(new EntryViewState(expandAction, false))
	}
	onEditClick(editAction: boolean) {
		this.setState(new EntryViewState(false, editAction))
	}
	onEditorCloseRequested() {
		this.setState(new EntryViewState(false, false))
	}
	render() {
		let en: WSEntryGetHTML = this.props.entry
		let entry: Entry = new Entry(en.EntryID, en.Title, "raw is missing", en.RawType, "tags is missing")
		if (this.state.editing) {
			return <div>
				<h2>Editing Entry: {en.Title}</h2>
				<EntryEditor entry={entry} editorCloseReq={e => this.onEditorCloseRequested()}/>
			</div>
		} else {
			if (this.state.expanded) {
				return <div>
					<h2 onClick={e => this.onExpandClick(false)}>{en.Title}
					</h2>
					<button onClick={e => this.onEditClick(true)}>Change</button>
					<div dangerouslySetInnerHTML={{__html: en.HTML}} />
				</div>
			} else {
				return <div><h2 onClick={e => this.onExpandClick(true)}>{en.Title}</h2></div>
			}
		}
	} // end of render function
} // end of class