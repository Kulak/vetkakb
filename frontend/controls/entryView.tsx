/*
entryView provides standard view of the entry from
a list of entries like search results.
*/
import React from 'react'
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
		expandInitially: boolean
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
		let fe = new WSFullEntry().initializeFromWSEntryGetHTML(props.entry)
		//console.log("EntryViewBox constructor, init state", fe)
		this.state = new EntryViewState(fe, props.expandInitially, false, false)
		User.Current()
		.then(function(json) {
			let user = json as WSUserGet
			let canEdit = user.Clearances == 8
			this.setState(new EntryViewState(fe, this.state.expanded, this.state.editing, canEdit))
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
				WSFullEntry.fromData(jsonEntry as WSFullEntry)
				.then((entry: WSFullEntry) => {
					console.log("entryView jsonEntry", jsonEntry)
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
		//console.log("entryView: render entry", fe)
		if (this.state.editing) {
			// in editing state
			return <EntryEditor entry={fe} editorCloseReq={fe => this.onEditorCloseRequested(fe)} />
		} else {
			// viewing; not editing
			if (this.state.expanded) {
				// is expanded
				let icon: JSX.Element = (null)
				if (fe.TitleIcon.length > 0) {
					icon = (<img className="uk-card uk-card-default uk-float-left" src={fe.TitleIcon} />)
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
				return (<article className="uk-article">
					<div className='uk-card uk-card-default'>
						<h1 className="uk-card-title uk-card-title-hover"
								onClick={e => this.onExpandClick(e, false)}>{fe.Title}
							{editButton}
						</h1>
					{icon}
					{entryBody}
					<a href={fe.permalink()}>Permalink</a>
					</div>
				</article>)
			} else {
				// not expanded
				let icon: JSX.Element = (null)
				if (fe.TitleIcon.length > 0) {
					icon = (<img className="uk-card uk-float-left uk-card-small" src={fe.TitleIcon} />)
				}
				return (
					<div className="uk-card uk-card-default">
						<h1 className="uk-card-title uk-card-title-hover" onClick={e => this.onExpandClick(e, true)}><a href={fe.permalink()}>{fe.Title}</a></h1>
						<div className='uk-card-body'>
							{icon}
							<p className="">{fe.Intro}</p>
						</div>
					</div>)
			}
		}
	} // end of render function
} // end of class
