
import React from 'react';
import * as ee from '../controls/entryCreate';
import {Recent} from '../controls/recent';
import {Search} from '../controls/search';
import {WSFullEntry} from '../model/wsentry'
import {EntryViewBox} from '../controls/entryView'
import {UserBox} from '../controls/userBox'

declare var ThemeName: string;
declare var ThemeURL: string;

class SingleEntryPageState {
  constructor(
  ) {}
}

export class SingleEntryPage extends React.Component<Object, SingleEntryPageState> {

  public constructor(props: Object, context) {
    super(props, context)
    this.state = new SingleEntryPageState()
  }

	render() {
    return <p>Single Entry Page</p>
  } // end of render function

} // end of class
