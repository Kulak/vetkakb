/*
userBox provides logon link or displays user name
*/

import * as React from 'react'
import {WSUserGet} from '../model/wsuser'
import {User} from '../common/user'

declare var SiteID: number;
declare var ZonePath: string

export interface UserBoxProps {
}

class UserBoxState {
	constructor(
		public user: WSUserGet
	) {}
}

export class UserBox extends React.Component<UserBoxProps, UserBoxState> {
	public constructor(props: UserBoxProps, context) {
		super(props, context)
		this.state = new UserBoxState(null)
		User.Current()
		.then(function(json) {
			let user = json as WSUserGet
			this.setState(new UserBoxState(user))
		}.bind(this))
		.catch(function(err) {
			console.log("error getting session user: ", err)
		}.bind(this))
	}
	render() {
		let u = this.state.user
		if (u != null && u.NickName != 'Guest') {
			return <div><img src={u.AvatarURL} className="avatar" />{u.NickName}</div>
		} else {
			return <form action={ZonePath + '/api/auth'}>
				<input type='hidden' name='state' value={SiteID} />
				<button className='leftStack' name='provider' value='gplus'>Login</button>
			</form>
		}
	} // end of render function
} // end of class