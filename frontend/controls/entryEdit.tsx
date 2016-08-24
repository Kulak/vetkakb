import * as React from 'react'
import {Entry} from '../model/entry'

export class EntryProps {
	constructor(public id: number){}
}

class EntryState {
	constructor(public isEditing: boolean = false,
		public entry: Entry = new Entry()) {}
}

export class EntryBox extends React.Component<EntryProps, EntryState> {
	public constructor(props: EntryProps, context) {
		super(props, context)
		this.state = new EntryState(false, new Entry());
	}
	onEditClick() {
		this.setState(new EntryState(true, this.state.entry))
	}
	onEditCancelClick() {
		this.setState(new EntryState(false, this.state.entry))
	}
	onEditSaveClick() {
		// save data
		this.setState(new EntryState(false, this.state.entry))
	}
	onEntryTitleChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EntryState(), this.state) as EntryState;
		state.entry.title = (event.target as any).value
		this.setState(state)
	}
	onEntryOrigBodyChange(event: React.FormEvent) {
		let state = (Object as any).assign(new EntryState(), this.state) as EntryState;
		state.entry.origBody = (event.target as any).value
		this.setState(state)
	}
	render() {
		if (this.state.isEditing) {
			return <div>
				<p>
					<button onClick={e => this.onEditCancelClick()}>Cancel Changes</button>
					<button onClick={e => this.onEditSaveClick()}>Save and Close</button>
				</p>
				<p>
					<label>Title:</label><input type="text" value={this.state.entry.title} onChange={e => this.onEntryTitleChange(e)} />
				</p>
				<p>
					<label>Raw Text:</label><br />
					<textarea value={this.state.entry.origBody} onChange={e => this.onEntryOrigBodyChange(e)} />
				</p>
				<label>Preview:</label>
				<pre>{this.state.entry.origBody}</pre>
			</div>
		} else {
			return <div><button onClick={e => this.onEditClick()}>New Entry</button></div>
		}
	} // end of render function
} // end of class