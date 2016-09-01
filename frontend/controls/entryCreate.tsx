/*
entryEdit provides a way to create new and update existing
entries.
*/

import * as React from 'react'
//import {Entry} from '../model/entry'
import {WSFullEntry, RawType} from '../model/wsentry'
import {EntryEditor, EditorProps} from './entryEditor'

export interface EntryCreateProps {
}

class EntryCreateState {
	constructor(
		public isEditing: boolean = false,
		public entry: WSFullEntry = new WSFullEntry()
	) {}
}

/*
	EntryBox provides basic editor to create or update an entry
	and save to the server.
*/
export class EntryCreateBox extends React.Component<EntryCreateProps, EntryCreateState> {
	public constructor(props: EntryCreateProps, context) {
		super(props, context)
		this.state = new EntryCreateState()
		this.state.entry.RawType = RawType.PlainText
	}
	onEditClick() {
		this.setState(new EntryCreateState(true, this.state.entry))
	}
	onEditorCloseRequested() {
		this.setState(new EntryCreateState(false, this.state.entry))
	}
	render() {
		if (this.state.isEditing) {
			return <div>
				<EntryEditor entry={this.state.entry} editorCloseReq={e => this.onEditorCloseRequested()} />
			</div>
		} else {
			return <div><button onClick={e => this.onEditClick()}>New Entry</button></div>
		}
	} // end of render function
} // end of class