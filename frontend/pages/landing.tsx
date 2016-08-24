
import * as React from 'react';
import * as ee from '../controls/entryEdit';

export class LandingPage extends React.Component<Object, Object> {
  public constructor(props: Object, context) {
    super(props, context)
  }
	render() {
      return <div>
        <h1>Landing page</h1>
        <ee.EntryBox />
      </div>
  } // end of render function
} // end of class