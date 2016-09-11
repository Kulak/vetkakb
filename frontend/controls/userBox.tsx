/*
userBox provides logon link or displays user name
*/

import * as React from 'react'
import {WSUserGet} from '../model/wsuser'
import {DataService} from '../common/dataService'

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
		DataService.get("/api/session/user")
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
		if (u != null) {
			return <div><img src={u.AvatarURL} className="avatar" />{u.Name}</div>
		} else {
			return <form action='api/auth'>
				<button className='leftStack' name='provider' value='gplus'>Login</button>
			</form>
		}
	} // end of render function
} // end of class