
import * as React from 'react'
import * as ReactDOM from 'react-dom'
import {LandingPage} from './pages/landing'
import {SingleEntryPage} from './pages/singleEntry'
import { Router, Route, Link, browserHistory } from 'react-router'
import {Recent} from './controls/recent'
import {Search} from './controls/search'
import {EntryCreateBox} from './controls/entryCreate'
import {UserBox} from './controls/userBox'
import {WSFullEntry} from './model/wsentry'
import {EntryLoaderBox} from './controls/entryLoader'
import {DataService} from './common/dataService'

declare var ThemeName: string
declare var ThemeURL: string
declare var ZonePath: string

console.log("ZonePath:", ZonePath)

let mainLayout = React.createClass({
	render: function() {
		return (<div>
      <h1>{document.title}</h1>
      <div className='toolbar'>
        <Link to={ZonePath + '/app/l/recent'} className='leftStack'>Recent</Link>
				<Link to={ZonePath + '/app/l/search'} className='leftStack'>Search</Link>
				<Link to={ZonePath + '/app/l/new'} className='leftStack'>New Entry</Link>
				<UserBox />
      </div>
      {this.props.children}
    </div>)
	}
})

let newEntry = React.createClass({

	onNewClose: function(fe: WSFullEntry) {
		if (fe.EntryID != 0) {
			browserHistory.push(ZonePath + '/app/l/e/' + fe.EntryID)
		} else {
			browserHistory.push(ZonePath + "/")
		}
	},

	render: function() {
		return (<EntryCreateBox editorCloseReq={fe => this.onNewClose(fe) }/>)
	}
})

ReactDOM.render((
	<Router history={browserHistory}>
		<Route path={ZonePath + "/index.html"} component={mainLayout} />
		<Route path={ZonePath + "/"} component={mainLayout} />
		<Route path={ZonePath + "/app/l"} component={mainLayout}>
			<Route path="recent" component={Recent} />
			<Route path="search" component={Search} />
			<Route path='new' component={newEntry} />
			<Route path="e/:entryID" component={EntryLoaderBox} />
			<Route path="e/:entryID/*" component={EntryLoaderBox} />
		</Route>
	</Router>
), document.getElementById('app'))
