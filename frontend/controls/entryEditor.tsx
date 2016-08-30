/*
EXPERIMENTAL: entryEditor focuses on changing content, saving it
*/

import * as React from 'react'
import {WSFullEntry} from '../model/wsentry'
import {WSEntryPut} from '../model/wsentry'
import {DataService} from '../common/dataService'

// example:
// declare type MyHandler = (myArgument: string) => void;
// var handler: MyHandler;


export type EditorCloseReqFunc = (e: any) => void;

export interface EditorProps extends React.Props<any>{
		entry: WSFullEntry;
		editorCloseReq: EditorCloseReqFunc;
}

class EditorState {
	constructor(
		public entry: WSFullEntry = new WSFullEntry()
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
		let pen: WSFullEntry = props.entry
		// make a copy of entry for easy cancellation
		this.state = new EditorState(new WSFullEntry(
			pen.EntryID, pen.Title, pen.Raw, pen.RawType, pen.Tags, pen.HTML, pen.Updated
		));
	}
	onEditCancelClick() {
		this.sendCloseRequest();
	}
	onEditSaveClick(close: boolean) {
		// save data
		let e: WSFullEntry = this.props.entry
		let base64: string = btoa(e.Raw)
		let wsEntry: WSEntryPut = new WSEntryPut(e.Title, base64, e.RawType, e.Tags)

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
		state.entry.Title = (event.target as any).value
		this.setState(new EditorState())
	}
	onEntryOrigBodyChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.Raw = (event.target as any).value
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
				<label>Title:</label><input type="text" value={this.state.entry.Title} onChange={e => this.onEntryTitleChange(e)} />
			</p>
			<p>
				<label>Raw Text:</label><br />
				<textarea value={this.state.entry.Raw} onChange={e => this.onEntryOrigBodyChange(e)} />
			</p>
			<label>Preview:</label>
			<pre>{this.state.entry.Raw}</pre>
		</div>
	} // end of render function
} // end of class