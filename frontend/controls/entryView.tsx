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
import {User} from '../common/user'
import {WSUserGet} from '../model/wsuser'

declare var ZonePath: string

export interface EntryViewProps {
		entry: WSEntryGetHTML
}

class EntryViewState {
	constructor(
		public fullEntry: WSFullEntry,
		public expanded: boolean,
		public editing: boolean,
		public canEdit: boolean,
	) {}
}

export class EntryViewBox extends React.Component<EntryViewProps, EntryViewState> {
	public constructor(props: EntryViewProps) {
		super(props)
		let pe = props.entry
		let fe = new WSFullEntry(pe.EntryID, pe.Title, pe.TitleIcon, null, pe.RawTypeName, "",
			pe.HTML, pe.Intro, pe.Updated)
		this.state = new EntryViewState(fe, false, false, false)
		User.Current()
		.then(function(json) {
			let user = json as WSUserGet
			let canEdit = user.Clearances == 8
			this.setState(new EntryViewState(fe, false, false, canEdit))
		}.bind(this))
		.catch(function(err) {
			console.log("error getting session user: ", err)
		}.bind(this))
	}
	onExpandClick(ev, expandAction: boolean) {
		ev.preventDefault()
		this.setState(new EntryViewState(this.state.fullEntry, expandAction, false, this.state.canEdit))
	}

	onEditClick(editAction: boolean) {
		if (editAction) {
			// load a full entry
			DataService.get(ZonePath + '/api/entry/' + this.props.entry.EntryID)
			.then((jsonEntry) => {
				//console.log("json text", jsonEntry)
				WSFullEntry.fromData(jsonEntry as WSFullEntry)
				.then((entry: WSFullEntry) => {
					this.setState(new EntryViewState(entry, false, editAction, this.state.canEdit))
				})
				.catch((err) => {
					// fromData does not reject
					console.error("Failed to obtain WSFullEntry from JSON response", err)
				})
			})
			.catch((err) => {
				console.log("err loading json: ", err)
			})
		} else {
			this.setState(new EntryViewState(this.state.fullEntry, false, editAction, this.state.canEdit))
		}
	}
	onEditorCloseRequested(fe: WSFullEntry) {
		console.log("entryView: set new EntryViewState", fe)
		this.setState(new EntryViewState(fe, true, false, this.state.canEdit))
	}
	render() {
		let fe: WSFullEntry = this.state.fullEntry
		console.log("entryView: render entry", fe)
		if (this.state.editing) {
			// in editing state
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
				let editButton = null
				if (this.state.canEdit) {
					editButton = <button onClick={e => this.onEditClick(true)}>Edit</button>
				}
				return <article className="uk-article">
					<div className='uk-panel uk-panel-box uk-panel-box-primary'>
						<h1 className="uk-article-title"
								onClick={e => this.onExpandClick(e, false)}>{fe.Title}
							{editButton}
						</h1>
					{icon}
					{entryBody}
					<a href={fe.permalink()}>Permalink</a>
					</div>
				</article>
			} else {
				// not expanded
				let icon: JSX.Element = (null)
				if (fe.TitleIcon.length > 0) {
					icon = (<img className="uk-thumbnail uk-float-left uk-thumbnail-mini" src={fe.TitleIcon} />)
				}
				return (
					<div className="uk-panel uk-panel-box uk-panel-box-primary uk-panel-box-primary-hove">
						<h1 className="uk-panel-title"><a href={fe.permalink()} onClick={e => this.onExpandClick(e, true)}>{fe.Title}</a></h1>
						{icon}
						<p className="">{fe.Intro}</p>
					</div>)
			}
		}
	} // end of render function
} // end of class
