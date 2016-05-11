import React, {PropTypes} from 'react';
import { withRouter } from 'react-router';
import AuthService from '../services/Auth';

class Login extends React.Component {

	static displayName = 'login';

	static propTypes = {
		location: PropTypes.object,
		router: PropTypes.object.isRequired,
	};

	constructor(props) {
		super(props);
		this.state = {
			loginStatus: {}
		};
	}

	componentDidMount() {
		if (AuthService.loggedIn) {
			this.props.router.replace('/');
			return;
		}
	}

	componentDidUpdate() {
		if (this.state.loginStatus && this.state.loginStatus.fulfilled) {
			const { location } = this.props;
			if (location.state && location.state.nextPathname) {
				this.props.router.replace(location.state.nextPathname);
			} else {
				this.props.router.replace('/');
			}
		}
	}

	submit(event) {
		event.preventDefault();
		this.setState({loginStatus: {pending: true}});
		AuthService.login(this.refs.username.value, this.refs.password.value).then(success => {
			this.setState({loginStatus: {pending: false, rejected: !success, fulfilled: success}});
		});
	}

	render() {

		let status = '';
		if (this.state.loginStatus.pending) {
			status = (<p>pending</p>);
		} else if (this.state.loginStatus.rejected) {
			status = (<p>rejected</p>);
		} else if (this.state.loginStatus.fulfilled) {
			status = (<p>fulfilled</p>);
		}

		return (
			<div>
				<form onSubmit={(event) => this.submit(event)}>
					<fieldset>
						<label htmlFor="username">Username</label>
						<input id="username" placeholder="username" ref="username" type="text" />
						<label htmlFor="password">Password</label>
						<input id="password" placeholder="password" ref="password" type="password" />
						<button className="button button-outline">login</button>
					</fieldset>
				</form>
				{status}
			</div>
		);

	}

}

export default withRouter(Login);
