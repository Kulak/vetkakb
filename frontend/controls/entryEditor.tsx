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

export type EditorCloseReqFunc = (fe: WSFullEntry) => void

export interface EditorProps extends React.Props<any>{
		entry: WSFullEntry
		editorCloseReq: EditorCloseReqFunc
}

class EditorState {
	constructor(
		public entry: WSFullEntry = new WSFullEntry(),
		public rawTypeName: string = ""
	) {}
}

/*
	EntryBox provides basic editor to create or update an entry
	and save to the server.
*/
export class EntryEditor extends React.Component<EditorProps, EditorState> {
	// used to manager refs
	// https://medium.com/@basarat/strongly-typed-refs-for-react-typescript-9a07419f807#.27st7hkss
	ctrls: {
		rawFile? :HTMLInputElement
	} = {}

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
		let r: Promise<any>
		if (this.state.rawTypeName.startsWith("Binary")) {
			r = this.onSaveBinary(close)
		} else {
			r = this.onSaveJson(close)
		}
		r.then((response) => {
			let fe: WSFullEntry = response
			fe.Raw = atob(fe.Raw)
			this.setState(new EditorState(fe, this.state.rawTypeName))
			if (close) {
				this.sendCloseRequest(fe);
			}
		})
		.catch((err) => {
			console.log("Failed to save", err)
			if (close) {
				this.sendCloseRequest(this.props.entry);
			}
		})

	}
	// saves binary file
	// example: https://www.raymondcamden.com/2016/05/10/uploading-multiple-files-at-once-with-fetch/
	onSaveBinary(close: boolean) :Promise<any> {
		// save data
		let e: WSFullEntry = this.state.entry
		let base64: string = btoa(e.Raw)
		let fileBlob = this.ctrls.rawFile.files[0]

		//var fd = new FormData();
		//fd.append('json', JSON.stringify(e))
		// first arg is the "formName" in multipartreader part
		//fd.append('rawFile', this.ctrls.rawFile.files[0]);

		return this.binaryContent(fileBlob)
		.then((fileContent) => {
			base64 = btoa(fileContent)
			console.log("file content: ", base64)

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
			return DataService.handleFetch("/entrybad", reqInit)
		})
		.catch((errMsg) => {
			console.log("file content error: ", errMsg)
		})

		//var fr = new FileReader()
		//fr.readAsDataURL(fileBlob)

	}

	binaryContent(fileBlob: Blob) :Promise<string> {
		return new Promise<string>((resolve, reject) => {
			let fr = new FileReader()
			fr.onload = (ev: Event) => {
				let data = (event.target as any).result as string;
				resolve(data)
			}
			fr.onerror = (ev: ErrorEvent) => {
				return reject(ev.error)
			}
			//fr.readAsDataURL(fileBlob)
			fr.readAsBinaryString(fileBlob)
		})
	}

	// saves standard scenario (JSON message)
	onSaveJson(close: boolean): Promise<any> {
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
		return DataService.handleFetch("/entry", reqInit)
	}
	// 	onSaveBinary(close: boolean) :Promise<WSFullEntry> {
	// 	return new Promise<WSFullEntry>((resolve, reject) => {
	// 		let reqInit: RequestInit
	// 		return DataService.handleFetch("/entry", reqInit)
	// 	})
	// }

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
	onEntryTagsChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.Tags = (event.target as any).value
		this.setState(state)
	}
	onRawTypeChange(rawType: number, name: string) {
		let state = (Object as any).assign(new EditorState(), this.state) as EditorState;
		state.entry.RawType = rawType
		state.rawTypeName = name
		this.setState(state)
	}
	render() {
		let rawPayload = <p>
				<label>Raw Text:</label><br />
				<textarea  value={this.state.entry.Raw} onChange={e => this.onEntryOrigBodyChange(e)} className='entryEdit' />
			</p>
		if (this.state.rawTypeName.startsWith("Binary")) {
			rawPayload = <p>
				<label>File upload:</label>
				<input type="file" ref={(input) => this.ctrls.rawFile = input} />
			</p>
		}
		return <div>
			<div className="toolbar entryHeader">
				<h2 className='leftStack'>Editing title:</h2>
				<input className='leftStack entryEdit' type="text" value={this.state.entry.Title}
					onChange={e => this.onEntryTitleChange(e)} />
			</div>
			<div className='toolbar'>
				<button className='leftStack' onClick={e => this.onEditSaveClick(false)}>Save</button>
				<RawTypeDropdown num={this.props.entry.RawType}
					rawTypeSelected={(num, name) => this.onRawTypeChange(num, name)} />
				<button className='leftStack' onClick={e => this.onEditSaveClick(true)}>OK</button>
				<button className='leftStack' onClick={e => this.onEditCancelClick()}>Cancel</button>
			</div>
			<p>
			</p>
			{rawPayload}
			<p>
				<label>Tags:</label><br />
				<input value={this.state.entry.Tags} onChange={e => this.onEntryTagsChange(e)} className='entryEdit' />
			</p>
		</div>
	} // end of render function
} // end of class