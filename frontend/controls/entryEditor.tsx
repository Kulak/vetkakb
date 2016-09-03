/*
entryEditor focuses on changing content, saving it
*/

import * as React from 'react'
import {WSFullEntry} from '../model/wsentry'
import {WSEntryPut, WSEntryPost, RawType} from '../model/wsentry'
import {DataService} from '../common/dataService'
import {RawTypeDropdown} from './rawTypeDropdown'

// example:
// declare type MyHandler = (myArgument: string) => void;
// var handler: MyHandler;

export type EditorCloseReqFunc = (fe: WSFullEntry) => void;

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
	sendCloseRequest(fe: WSFullEntry) {
		if (this.props.editorCloseReq != null) {
			this.props.editorCloseReq(fe)
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
		this.sendCloseRequest(this.props.entry);
	}
	onEditSaveClick(close: boolean) {
		// save data
		let e: WSFullEntry = this.state.entry
		let base64: string = btoa(e.Raw)
		let reqInit: RequestInit
		if (e.EntryID == 0) {
			// create new entry with PUT
			let wsEntry = new WSEntryPut(e.Title, base64, e.RawType, e.Tags)
			reqInit = DataService.newRequestInit("PUT", wsEntry)
		} else {
			// update existing entry with POST
			let wsEntry = new WSEntryPost(e.EntryID, e.Title, base64, e.RawType, e.Tags)
			reqInit = DataService.newRequestInit("POST", wsEntry)
		}
		DataService.handleFetch("/entry", reqInit)
		.then(function(jsonText) {
			console.log("json response", jsonText)
			let fe = jsonText as WSFullEntry
			this.setState(fe)
			if (close) {
				this.sendCloseRequest(fe);
			}
		}.bind(this))
		.catch(function(err) {
			console.log("response err: ", err)
			if (close) {
				this.sendCloseRequest(e)
			}
		}.bind(this))
	}
	onEntryTitleChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.Title = (event.target as any).value
		this.setState(state)
	}
	onEntryOrigBodyChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.Raw = (event.target as any).value
		this.setState(state)
	}
	onRawTypeChange(rawType: number) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.RawType = rawType
		this.setState(state)
	}
	render() {
		return <div>
			<p>
				<button onClick={e => this.onEditCancelClick()}>Cancel Changes</button>
				<button onClick={e => this.onEditSaveClick(false)}>Save and Edit</button>
				<button onClick={e => this.onEditSaveClick(true)}>Save and Close</button>
				<RawTypeDropdown num={this.props.entry.RawType}
					rawTypeSelected={e => this.onRawTypeChange(e)} />
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