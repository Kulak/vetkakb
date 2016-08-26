
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {LandingPage} from './pages/landing';

class AppProps {
	constructor(public defaultUrl: string){}
}

class AppState {
	constructor(public url: string) {}
}

class App extends React.Component<AppProps, AppState> {
	// context is currently {}
	constructor(props: AppProps, context) {
		super(props, context)
		this.state = new AppState(props.defaultUrl);
	}
	render() {
		return <LandingPage />
	}
} // end of App class

ReactDOM.render(
	<App defaultUrl="/" />,
	document.getElementById('app')
);
