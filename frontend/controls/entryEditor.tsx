/*
entryEditor focuses on changing content, saving it
*/

import * as React from 'react'
import {WSFullEntry} from '../model/wsentry'
import {WSEntryPost, RawType} from '../model/wsentry'
import {DataService} from '../common/dataService'
import {RawTypeDropdown} from './rawTypeDropdown'
import {WSRawType} from '../common/rawtypes'

declare var ZonePath: string


// example:
// declare type MyHandler = (myArgument: string) => void;
// var handler: MyHandler;

export type EditorCloseReqFunc = (fe: WSFullEntry) => void

export interface EditorProps extends React.Props<any>{
		entry: WSFullEntry
		editorCloseReq: EditorCloseReqFunc
}

interface EditorState {
		entry: WSFullEntry
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

	copyState(): EditorState {
		return Object.assign({}, this.state)
	}

	sendCloseRequest(fe: WSFullEntry) {
		if (this.props.editorCloseReq != null) {
			this.props.editorCloseReq(fe)
		}
	}
	public constructor(props: EditorProps, context) {
		super(props, context)
		// make a copy of entry for easy cancellation
		this.state = {entry: props.entry.copy()}
		//console.log("Edit is in state", this.state)
	}
	onEditCancelClick() {
		this.sendCloseRequest(this.props.entry);
	}
	onEditSaveClick(close: boolean) {
		this.onSaveBinary(close)
		.then((response) => {
			// get easy access to fields, but it is not a real WSFullEntry
			WSFullEntry.fromData(response as WSFullEntry)
			.then((fe) => {
				// the following line triggers React mesage:
				// 		react.js:20541 Warning: EntryEditor is changing a controlled input of type undefined to be uncontrolled.
				// 		Input elements should not switch from controlled to uncontrolled (or vice versa).
				// The setState call triggers React to think that state is different for controlled element.
				//    this.setState(new EditorState(fe, this.state.rawTypeName))
				// So, we simply rely on original state and merge most important properties here
				let nState = this.copyState()
				nState.entry.initializeFromWSFullEntry(fe)
				console.log("editor's setState on success (nState.entry, fe)", nState.entry, fe)
				this.setState(nState)

				if (close) {
					this.sendCloseRequest(fe);
				}
			})
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

		var fd = new FormData();

		let wsEntry = new WSEntryPost(e.EntryID, e.Title, e.TitleIcon, e.RawTypeName, e.Tags, e.Intro, e.Slug)
		fd.append('entry', JSON.stringify(wsEntry))

		console.log("on save entry (state, wsEntry)", this.state.entry, wsEntry)
		// check if raw is binary or some other text
		if (this.state.entry.RawTypeName.startsWith(WSRawType.Binary)) {
			let fileBlob = this.ctrls.rawFile.files[0]
			fd.append('rawFile', fileBlob);
		} else {
			fd.append('rawFile', e.Raw)
		}

		// POST in thise case PUTs new entry and updates existing
		let reqInit: RequestInit = DataService.newBareRequestInit("POST")
		reqInit.body = fd
		return DataService.handleFetch(ZonePath + "/binaryentry/", reqInit)
	}

	onEntryTitleChange(event: React.FormEvent) {
		let nState = this.copyState()
		nState.entry.Title = (event.target as any).value
		this.setState(nState)
	}
	onEntryTitleIconChange(event: React.FormEvent) {
		let nState = this.copyState()
		nState.entry.TitleIcon = (event.target as any).value
		this.setState(nState)
	}
	onEntryIntroChange(event: React.FormEvent) {
		let nState = this.copyState()
		nState.entry.Intro = (event.target as any).value
		this.setState(nState)
	}
	public onEntryOrigBodyChange: (event: React.FormEvent) => void = (event) => {
		let nState = this.copyState()
		nState.entry.Raw = (event.target as any).value
		this.setState(nState)
	}
	onEntrySlugChange(event: React.FormEvent) {
		let nState = this.copyState()
		nState.entry.Slug = (event.target as any).value
		this.setState(nState)
	}
	onEntryTagsChange(event: React.FormEvent) {
		let nState = this.copyState()
		nState.entry.Tags = (event.target as any).value
		this.setState(nState)
	}
	onRawTypeChange(rawTypeName: string) {
		let nState = this.copyState()
		nState.entry.RawTypeName = rawTypeName
		this.setState(nState)
	}

	render() {
		let rawPayload = <p>Data type name is not loaded yet.</p>
		if (this.state.entry.RawTypeName == "") {
			// do nothing
		} else if (this.state.entry.RawTypeName.startsWith(WSRawType.Binary)) {
			let image = <span />
			if (this.state.entry.RawTypeName == WSRawType.BinaryImage && this.state.entry.EntryID > 0) {
				image = <img className='' src={"re/" + this.state.entry.EntryID} />
			}
			// allow user to select new image
			rawPayload = <p>
				<label>File upload:</label>
				<input type="file" ref={(input) => this.ctrls.rawFile = input} />
				{image}
			</p>
		} else {
			rawPayload = <p>
				<label>Raw Text:</label><br />
				<textarea  value={this.state.entry.Raw} onChange={this.onEntryOrigBodyChange} className='entryEdit' />
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
				<RawTypeDropdown name={this.props.entry.RawTypeName}
					rawTypeSelected={(name) => this.onRawTypeChange(name)} />
				<button className='leftStack' onClick={e => this.onEditSaveClick(true)}>OK</button>
				<button className='leftStack' onClick={e => this.onEditCancelClick()}>Cancel</button>
			</div>
			{rawPayload}
			<p>
				<label>Slug:</label><br />
				<input value={this.state.entry.Slug} onChange={e => this.onEntrySlugChange(e)}
					className='entryEdit' />
			</p>
			<p>
				<label>Tags:</label><br />
				<input value={this.state.entry.Tags} onChange={e => this.onEntryTagsChange(e)}
					className='entryEdit' />
			</p>
			<p>
				<label>Title Icon URL:</label><br />
				<input value={this.state.entry.TitleIcon} onChange={e => this.onEntryTitleIconChange(e)}
					className='entryEdit' />
			</p>
			<p>
				<label>Introduction:</label><br />
				<textarea  value={this.state.entry.Intro} onChange={e => this.onEntryIntroChange(e)}
					className='entryEdit entryEdit-short' />
			</p>
		</div>
	} // end of render function
} // end of class


// We cannot simply display an image as we are not certain of its MIME Type.
// The code below is for reference purpose.
//
// if (this.state.rawTypeName == WSRawType.BinaryImage && this.state.entry.Raw != null) {
// 	// display existing image
// 	let imgUrl = "data:image/png;base64," + btoa(this.state.entry.Raw)
// 	rawPayload = <p>Image:
// 		<img src={imgUrl} />
// 	</p>
// }
