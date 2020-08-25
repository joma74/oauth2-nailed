import React from "react"
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom"
import "./App.css"
import AuthCodeRedirect from "./AuthCodeRedirect"

class App extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			access_token: "",
			expires_in: "",
			session_state: "",
			token_type: "",
		}
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

	render() {
		return (
			<Router>
				<div>
					<nav>
						<div className="App">
							<h1>Implicit Grant Type</h1>
						</div>
						<ul>
							<li>
								<Link to="/">Home</Link>
							</li>
							<li>
								<Link to="/login">Login</Link>
							</li>
							<li>
								<Link to="/services">Services</Link>
							</li>
							<li>
								<Link to="/logout">Logout</Link>
							</li>
						</ul>
					</nav>
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
							<Logout />
						</Route>
						<Route path="/">
							<Home />
						</Route>
					</Switch>
				</div>
			</Router>
		)
	}
}

function Login() {
	window.location.href =
		"http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/auth?client_id=oauth-nailed-app-2-implicit-grant&response_type=token&redirect_uri=http://localhost:3000/authCodeRedirect&scope=billingService"
	return null
}

function Services({ _accessToken }) {
	const formData = new FormData()
	formData.append("access_token", _accessToken)
	fetch("http://localhost:9111/billing/v1/services", {
		method: "POST",
		body: formData,
	})
		.then((rs) => rs.json())
		.then((data) => {
			console.log(data)
		})
	return <h2>Services</h2>
}

function Logout() {
	return <h2>Logout</h2>
}

function Home() {
	return <h2>Home</h2>
}

export default App
