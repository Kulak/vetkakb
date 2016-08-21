
/// <reference path="./typings/modules/react/index.d.ts" />

import * as React from 'react';
import * as ReactDOM from 'react-dom';

class App extends React.Component<any> {
  constructor() {
  }
  render() {
    return <div>Hello world!</div>
  }
}

ReactDOM.render(
  <App />,
  document.getElementById('app')
);
