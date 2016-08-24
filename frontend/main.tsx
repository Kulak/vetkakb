
//  / <reference path="./typings/modules/react/index.d.ts" />

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
    //new CommonData()
    new LandingPage({}, {})
    //return <LandingPage />
    // if (this.state.url == '/') {
       return <LandingPage />
    // } else {
    //   return <div>Unknown url: {this.state.url}</div>
    // }
  }
} // end of App class

ReactDOM.render(
  <App defaultUrl="/" />,
  document.getElementById('app')
);
