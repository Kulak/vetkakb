/*
A dropdown control to select entry type.
*/

import * as React from 'react'
import {WSRawType} from '../common/rawtypes'

export type RawTypeSelectedFunc = (num: number, name: string) => void;

export interface RawTypeDropdownProps extends React.Props<any> {
	num: number,
	rawTypeSelected: RawTypeSelectedFunc
}

class RawTypeDropdownState {
	constructor(
		public num: number = 0,
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
				let s = new RawTypeDropdownState(this.props.num, rawTypes)
				this.setState(s)
			}.bind(this))
			.catch(function(err) {
				console.log("WSRawType err: ", err)
			})
	}

	onSelectionChange(e) {
		console.log("RawTypeDD on change: ", e)
		if (e.target.selectedOptions.length) {
			let num = e.target.selectedOptions[0].value  // or label
			this.setState(new RawTypeDropdownState(num, this.state.rawTypes))
			console.log("RawTypeDD selected : ", num)
			let name: string = this.state.rawTypes.find(function(each) {
				return each.TypeNum == num
			}).Name
			this.props.rawTypeSelected(parseInt(num), name)
		}
	}

	render() {
		let rawTypes = <span>Loading raw type...</span>
		if (this.state.rawTypes != null) {
			let options = this.state.rawTypes.map(function(each) {
				return <option key={each.TypeNum}
					value={each.TypeNum}>{each.Name}</option>
			})
			rawTypes = <select value={this.state.num}
				onChange={e => this.onSelectionChange(e)}>
							{options}
            </select>
		}
		return rawTypes
	}

}