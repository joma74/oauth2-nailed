import React, { Component } from "react"
import {
	BrowserRouter as Router,
	Switch,
	Route,
	Link,
	Redirect,
} from "react-router-dom"
import "./App.css"
import AuthCodeRedirect from "./AuthCodeRedirect"
import Services from "./Services"

class App extends React.Component {
	constructor(props) {
		super(props)
		this.initialState = {
			access_token: "???",
			expires_in: "???",
			session_state: "???",
			token_type: "???",
		}
		this.state = Object.assign({}, this.initialState)
	}

	/**
	 *
	 * @param {string} k
	 * @param {string} v
	 */
	_setStateValue = (k, v) => {
		if (this.state[k] !== v) {
			this.setState({ [k]: v })
		}
	}

	_setState = (newState) => {
		this.setState(Object.assign(this.state, newState))
	}

	render() {
		return (
			<Router>
				<h1>
					<Link to="/">oauth-nailed-app-2-implicit-grant</Link>
				</h1>
				<Link to="/login">
					<button>Login with Keycloak</button>
				</Link>
				<Link to="/services">
					<button>Show billing services</button>
				</Link>
				<Link to="/logout">
					<button>Logout from Keycloak</button>
				</Link>
				<div>
					<h2>Infos from Keycloak</h2>
					<p>
						session state:<pre>{this.state.session_state}</pre>
					</p>
					<p>
						access token:<pre>{this.state.access_token}</pre>
					</p>
					<p>
						expires_in:<pre>{this.state.expires_in}</pre>
					</p>
					<p>
						token type:<pre>{this.state.token_type}</pre>
					</p>
				</div>
				{/* A <Switch> looks through its children <Route>s and
                renders the first one that matches the current URL. */}
				<Switch>
					<Route path="/login">
						<Login />
					</Route>
					<Route path="/services">
						<Services _accessToken={this.state.access_token} />
					</Route>
					<Route
						path="/authCodeRedirect"
						render={() => (
							<AuthCodeRedirect _setStateValue={this._setStateValue} />
						)}
					></Route>
					<Route path="/logout">
						<Logout
							_setState={this._setState}
							_initialState={this.initialState}
						/>
					</Route>
					<Route path="/">
						<Home />
					</Route>
				</Switch>
			</Router>
		)
	}
}

function Login() {
	window.location.href =
		"http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/auth?client_id=oauth-nailed-app-2-implicit-grant&response_type=token&redirect_uri=http://localhost:9109/authCodeRedirect&scope=billingService"
	return null
}

class Logout extends Component {
	constructor(props) {
		super(props)
		this.state = {}
	}

	componentDidMount() {
		const { _setState, _initialState } = this.props
		_setState(_initialState)
		return null
	}

	render() {
		return <Redirect to="/" />
	}
}

function Home() {
	return null
}

export default App
