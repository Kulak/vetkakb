/*
Search is a top level control that combines search conditions
with action to load data based on selected conditions
and with display of search results.

It is a core search controller.
*/

import * as React from 'react'

export class SearchProps {}

class SearchState {}

export class Search extends React.Component<SearchProps, SearchState> {
	onRecentClick() {
		console.log("in onRecentClick")
	}
	public constructor(props: SearchProps, context: any) {
		super(props, context)
	}  // end of constructor
	render() {
		return <div>
			<button onClick={e => this.onRecentClick()}>Recent</button>
		</div>
	}
}  // end of class