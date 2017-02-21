/*
entryEdit provides a way to create new and update existing
entries.
*/

import React from 'react'
//import {Entry} from '../model/entry'
import {WSFullEntry, RawType} from '../model/wsentry'
import {EntryEditor, EditorProps, EditorCloseReqFunc} from './entryEditor'

export interface EntryCreateProps {
		editorCloseReq: EditorCloseReqFunc
}

class EntryCreateState {
	constructor(
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
		this.state.entry.RawTypeName = RawType.PlainText
	}
	onEditorCloseRequested(fe: WSFullEntry) {
		this.setState(new EntryCreateState(this.state.entry))
		this.props.editorCloseReq(fe)
	}
	render() {
		return <div>
			<EntryEditor entry={this.state.entry} editorCloseReq={fe => this.onEditorCloseRequested(fe)} />
		</div>
	} // end of render function
} // end of class