/*
Search is a top level control that combines search conditions
with action to load data based on selected conditions
and with display of search results.

It is a core search controller.
*/

import * as React from 'react'

export class SearchProps {}

class SearchState {}

export class Search extends React.Component<SearchProps, Search> {
	public constructor(props: SearchProps, context: any) {
		super(SearchProps, context)
	}  // end of constructor
	render() {
		return <div>Search Bar from menu bar.</div>
	}
}  // end of class