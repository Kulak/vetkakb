
import * as React from 'react';
import * as ee from '../controls/entryCreate';
import {Recent} from '../controls/recent';
import {WSFullEntry} from '../model/wsentry'

class LandingPageState {
  constructor(
    public path: string
  ) {}
}

export class LandingPage extends React.Component<Object, LandingPageState> {

  public constructor(props: Object, context) {
    super(props, context)
    this.state = new LandingPageState('')
  }

  onNewClick() {
    this.setState(new LandingPageState('new'))
  }

  onRecentClick() {
    this.setState(new LandingPageState('recent'))
  }

  onNewClose(fe: WSFullEntry) {
    this.setState(new LandingPageState(''))
  }

	render() {
    let body = <div />
    if (this.state.path == 'new') {
      body = <ee.EntryCreateBox editorCloseReq={fe => this.onNewClose(fe) } />
    } else if (this.state.path == 'recent') {
      body = <Recent />
    }

    return <div>
      <h1>Vetka Information Management System</h1>
      <div className='toolbar'>
        <button className='leftStack' onClick={e => this.onNewClick()} >New Entry</button>
        <button className='leftStack' onClick={e => this.onRecentClick()} >Recent Entires</button>
      </div>
      {body}
    </div>
  } // end of render function

} // end of class