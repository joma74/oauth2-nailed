import React, { Component } from "react"
import { withRouter, Redirect } from "react-router-dom"

class AuthCodeRedirect extends Component {
	constructor(props) {
		super(props)
		this.state = {}
	}

	static getDerivedStateFromProps({ _setStateValue, location }, state) {
		// get access token
		const hashURL = location.hash
		const hashURLAsMap = hashURL
			.substr(1)
			.split("&")
			.reduce((acc, item) => {
				const kv = item.split("=")
				const k = kv[0]
				const v = kv[1]
				acc[k] = v
				return acc
			}, {})
		_setStateValue("access_token", hashURLAsMap["access_token"])
		_setStateValue("expires_in", hashURLAsMap["expires_in"])
		_setStateValue("session_state", hashURLAsMap["session_state"])
		_setStateValue("token_type", hashURLAsMap["token_type"])
		return null
	}
	render() {
		return <Redirect to="/" />
	}
}

export default withRouter(AuthCodeRedirect)
