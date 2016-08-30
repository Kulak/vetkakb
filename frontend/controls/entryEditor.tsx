/*
EXPERIMENTAL: entryEditor focuses on changing content, saving it
*/

import * as React from 'react'
import {Entry} from '../model/entry'
import {WSEntryPut} from '../model/wsentry'
import {DataService} from '../common/dataService'

// example:
// declare type MyHandler = (myArgument: string) => void;
// var handler: MyHandler;


export type EditorCloseReqFunc = (e: any) => void;

export interface EditorProps extends React.Props<any>{
		entry: Entry;
		editorCloseReq: EditorCloseReqFunc;
}

class EditorState {
	constructor(
		public entry: Entry = new Entry()
	) {}
}

/*
	EntryBox provides basic editor to create or update an entry
	and save to the server.
*/
export class EntryEditor extends React.Component<EditorProps, EditorState> {
	sendCloseRequest() {
		if (this.props.editorCloseReq != null) {
			this.props.editorCloseReq({})
		}
	}
	public constructor(props: EditorProps, context) {
		super(props, context)
		let pen: Entry = props.entry
		// make a copy of entry for easy cancellation
		this.state = new EditorState(new Entry(
			pen.entryID, pen.title, pen.raw, pen.rawType, pen.tags
		));
	}
	onEditCancelClick() {
		this.sendCloseRequest();
	}
	onEditSaveClick(close: boolean) {
		// save data
		let e: Entry = this.props.entry
		let base64: string = btoa(e.raw)
		let wsEntry: WSEntryPut = new WSEntryPut(e.title, base64, e.rawType, e.tags)

		DataService.put('/entry/', wsEntry)
		.then(function(jsonText) {
			console.log("PUT json response", jsonText)
		})
		.catch(function(err) {
			console.log("PUT err: ", err)
		})

		if (close) {
			this.sendCloseRequest();
		}
	}
	onEntryTitleChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.title = (event.target as any).value
		this.setState(new EditorState())
	}
	onEntryOrigBodyChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.raw = (event.target as any).value
		this.setState(state)
	}
	render() {
		return <div>
			<p>
				<button onClick={e => this.onEditCancelClick()}>Cancel Changes</button>
				<button onClick={e => this.onEditSaveClick(false)}>Save and Edit</button>
				<button onClick={e => this.onEditSaveClick(true)}>Save and Close</button>
			</p>
			<p>
				<label>Title:</label><input type="text" value={this.state.entry.title} onChange={e => this.onEntryTitleChange(e)} />
			</p>
			<p>
				<label>Raw Text:</label><br />
				<textarea value={this.state.entry.raw} onChange={e => this.onEntryOrigBodyChange(e)} />
			</p>
			<label>Preview:</label>
			<pre>{this.state.entry.raw}</pre>
		</div>
	} // end of render function
} // end of class