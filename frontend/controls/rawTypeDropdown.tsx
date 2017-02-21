import React from 'react'
import {WSRawType} from '../common/rawtypes'

/*
A dropdown control to select entry type.
*/

export type RawTypeSelectedFunc = (rawTypeName: string) => void;

export interface RawTypeDropdownProps extends React.Props<any> {
	name: string,
	rawTypeSelected: RawTypeSelectedFunc
}

class RawTypeDropdownState {
	constructor(
		public name: string = "",
		public rawTypes: Array<WSRawType> = null
	) {}
}

export class RawTypeDropdown extends React.Component<RawTypeDropdownProps, RawTypeDropdownState> {

	public constructor(props: RawTypeDropdownProps) {
		super(props)
		this.state = new RawTypeDropdownState()
		WSRawType.List()
			.then(function(rawTypes: Array<WSRawType>) {
				console.log("WSRawType LIST", rawTypes)
				let s = new RawTypeDropdownState(this.props.name, rawTypes)
				this.setState(s)
				// send initial notification of raw type name
				this.props.rawTypeSelected(props.name)
			}.bind(this))
			.catch(function(err) {
				console.log("WSRawType err: ", err)
			})
	}

	onSelectionChange(e) {
		console.log("RawTypeDD on change: ", e)
		if (e.target.selectedOptions.length) {
			let name = e.target.selectedOptions[0].value  // or label
			this.setState(new RawTypeDropdownState(name, this.state.rawTypes))
			console.log("RawTypeDD selected : ", name)
			this.props.rawTypeSelected(name)
		}
	}

	render() {
		let rawTypes = <span>Loading raw type...</span>
		if (this.state.rawTypes != null) {
			let options = this.state.rawTypes.map(function(each) {
				return <option key={each.Name}
					value={each.Name}>{each.Name}</option>
			})
			rawTypes = <select value={this.state.name}
				onChange={e => this.onSelectionChange(e)}>
							{options}
            </select>
		}
		return rawTypes
	}

}