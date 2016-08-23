
//  / <reference path="./typings/modules/react/index.d.ts" />

import * as React from 'react';
import * as ReactDOM from 'react-dom';

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
    if (this.state.url == '/') {
      return <div>Default page</div>
    } else {
      return <div>Unknown url: {this.state.url}</div>
    }
  }
} // end of App class

ReactDOM.render(
  <App defaultUrl="/" />,
  document.getElementById('app')
);
