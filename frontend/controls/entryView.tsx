/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import * as React from 'react'
import {WSEntryGetHTML, WSFullEntry} from '../model/wsentry'
import {EntryEditor, EditorProps} from './entryEditor'
import {Entry} from '../model/entry'
import {DataService} from '../common/dataService'
import {WSRawType} from '../common/rawtypes'

declare var ZonePath: string

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
		let fe = new WSFullEntry(pe.EntryID, pe.Title, pe.TitleIcon, null, pe.RawTypeName, "",
			pe.HTML, pe.Intro, pe.Updated)
		this.state = new EntryViewState(fe, false, false);
	}
	onExpandClick(expandAction: boolean) {
		this.setState(new EntryViewState(this.state.fullEntry, expandAction, false))
	}

	b64toBlob = (b64Data, contentType='', sliceSize=512) => {
		const byteCharacters = atob(b64Data);
		const byteArrays = [];

		for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
			const slice = byteCharacters.slice(offset, offset + sliceSize);

			const byteNumbers = new Array(slice.length);
			for (let i = 0; i < slice.length; i++) {
				byteNumbers[i] = slice.charCodeAt(i);
			}

			const byteArray = new Uint8Array(byteNumbers);

			byteArrays.push(byteArray);
		}

		const blob = new Blob(byteArrays, {type: contentType});
		return blob;
	}

	onEditClick(editAction: boolean) {
		if (editAction) {
			// load a full entry
			DataService.get(ZonePath + '/api/entry/' + this.props.entry.EntryID)
			.then(function(jsonEntry) {
				console.log("json text", jsonEntry)
				let entry = jsonEntry as WSFullEntry
				// don't convert null, because it atob(null) returns "ée"
				if (entry.Raw != null) {
					let blob = this.b64toBlob(entry.Raw)
					let reader = new FileReader()
					reader.addEventListener('loadend', function() {
						// listener is called when readAsText is completed; promise could be used here
						// For ISO-8859-1 there's no further conversion required
						console.log("FIXED:", reader.result)
						entry.Raw = reader.result
						this.setState(new EntryViewState(entry, false, editAction))
					}.bind(this))
					reader.readAsText(blob)
				} else {
					entry.Raw = ""
					this.setState(new EntryViewState(entry, false, editAction))
				}
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
		//console.log("entryView: render entry", fe)
		if (this.state.editing) {
			return <EntryEditor entry={fe} editorCloseReq={fe => this.onEditorCloseRequested(fe)} />
		} else {
			// viewing; not editing
			if (this.state.expanded) {
				// is expanded
				let icon: JSX.Element = (null)
				if (fe.TitleIcon.length > 0) {
					icon = (<img className="uk-thumbnail uk-float-left" src={fe.TitleIcon} />)
				}
				let entryBody
				if (fe.RawTypeName == WSRawType.BinaryImage) {
					entryBody = <img className='' src={"re/" + fe.EntryID} />
				} else {
					entryBody = <div dangerouslySetInnerHTML={{__html: fe.HTML}} />
				}
				return <article className="uk-article">
					<nav className='uk-navbar'>
						<ul className="uk-navbar-nav">
							<li><h1 className="uk-article-title"
								onClick={e => this.onExpandClick(false)}>{fe.Title}</h1></li>
							<li><button onClick={e => this.onEditClick(true)}>Edit</button></li>
						</ul>
					</nav>
					{icon}
					{entryBody}
					<hr className="uk-article-divider" />
				</article>
			} else {
				// not expanded
				let icon: JSX.Element = (null)
				if (fe.TitleIcon.length > 0) {
					icon = (<img className="uk-thumbnail uk-float-left uk-thumbnail-mini" src={fe.TitleIcon} />)
				}
				return (
					<div className="uk-panel uk-clearfix">
						<h1 className="uk-panel-title"> <a onClick={e => this.onExpandClick(true)}>{fe.Title}</a></h1>
						{icon}
						<p className="">{fe.Intro}</p>
					</div>)
			}
		}
	} // end of render function
} // end of class
