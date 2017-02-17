
import * as React from 'react'
import * as ReactDOM from 'react-dom'
import {LandingPage} from './pages/landing'
import {SingleEntryPage} from './pages/singleEntry'
import { Router, Route, IndexRoute, Link, browserHistory, Redirect } from 'react-router'
import {Recent} from './controls/recent'
import {Search} from './controls/search'
import {EntryCreateBox} from './controls/entryCreate'
import {UserBox} from './controls/userBox'
import {WSFullEntry} from './model/wsentry'
import {EntryLoaderBox} from './controls/entryLoader'
import {SlugLoaderBox} from './controls/slugLoader'
import {Dashboard} from './controls/dashboard'
import {DataService} from './common/dataService'

declare var ThemeName: string
declare var ThemeURL: string
declare var ZonePath: string

console.log("ZonePath:", ZonePath)

let mainLayout = React.createClass({
	componentDidMount: function() {
		/* this is a workaround for custom attributes in React
			https://jsfiddle.net/peterjmag/kysymow0/
			The following line changes
			<nav ref='navbarContainer' className='uk-navbar-container'>
			to
			<nav ref='navbarContainer' className='uk-navbar-container' uk-navbar >

			this.refs.navbarContainer.getDOMNode().setAttribute('uk-navbar', 'true')

			it does not seem to work
		*/
	},

	render: function() {
		return (<div className='uk-container uk-width-2-3'>
			{/* center logo image; width is set to the size of the image */}
			<div style="position: relative;margin-left: auto;margin-right: auto;width: 300px;">
				<img src={ZonePath + '/res/logo.png'}/>
			</div>
      <nav className='uk-navbar-container'>
				<div className='uk-navbar-left'>
					<ul className="uk-navbar-nav">
						<li><Link to={ZonePath + '/app/recent'} className=''>Recent</Link></li>
						<li><Link to={ZonePath + '/app/search'} className=''>Search</Link></li>
						<li><Link to={ZonePath + '/app/new'} className=''>New Entry</Link></li>
					</ul>
					<UserBox />
				</div>
      </nav>
      {this.props.children}
    </div>)
	}
})

let newEntry = React.createClass({

	onNewClose: function(fe: WSFullEntry) {
		if (fe.EntryID != 0) {
			browserHistory.push(ZonePath + '/app/e/' + fe.EntryID)
		} else {
			browserHistory.push(ZonePath + "/")
		}
	},

	render: function() {
		return (<EntryCreateBox editorCloseReq={fe => this.onNewClose(fe) }/>)
	}
})

//
ReactDOM.render((
	<Router history={browserHistory}>
		{/* Redirect from / to /app/ */}
		{/* <Redirect from={ZonePath + "/"} to={ZonePath + "/app"} /> */}
		<Route path={ZonePath + "/"} component={mainLayout} >
			<IndexRoute component={Dashboard} />
			<Route path={"s/:slug"} component={SlugLoaderBox} />
		</Route>
		<Route path={ZonePath + "/app"} component={mainLayout}>
			<IndexRoute component={Dashboard} />
			<Route path="recent" component={Recent} />
			<Route path="search" component={Search} />
			<Route path='new' component={newEntry} />
			<Route path="e/:entryID" component={EntryLoaderBox} />
			<Route path="e/:entryID/*" component={EntryLoaderBox} />
		</Route>
	</Router>
), document.getElementById('app'))
