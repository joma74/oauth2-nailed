import React, { Component } from "react"

class Services extends Component {
	constructor(props) {
		super(props)
		this.serviceEndpointURL = "http://localhost:9111/billing/v1/services"
		this.initialState = {
			services: [
				"😞 Billing services not available or access denied, check the log for cause",
			],
		}
		this.state = Object.assign({}, this.initialState)
	}

	componentDidMount() {
		const { _accessToken } = this.props
		const formData = new FormData()
		formData.append("access_token", _accessToken)
		fetch(this.serviceEndpointURL, {
			method: "POST",
			body: formData,
		})
			.then((rs) => {
				if (!rs.ok) {
					throw new Error(
						`Expected service response to be of status 200 but was ${rs.status}`,
					)
				}
				return rs.json()
			})
			.then((services) => {
				this.setState(Object.assign({}, services))
			})
			.catch((error) => {
				console.error(error)
				this.setState(Object.assign({}, this.initialServices))
			})
	}

	render() {
		const items = this.state.services.map((service, index) => (
			<li key={index}>{service}</li>
		))
		return (
			<div>
				<h2>List of billing services (from the Protected Resource)</h2>
				<ul>{items}</ul>
			</div>
		)
	}
}

export default Services
