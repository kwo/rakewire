import React, {PropTypes} from 'react';
import { withRouter } from 'react-router';
import AuthService from '../services/Auth';

import {Card, CardActions, CardTitle, CardText} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import TextField from 'material-ui/TextField';

const style = {
	card: {
	},
	cardheader: {
	},
	field: {
	},
	input: {
	}
};

class Login extends React.Component {

	static displayName = 'login';

	static contextTypes = {
		router: PropTypes.object.isRequired
	}

	static propTypes = {
		location: PropTypes.object
	};

	constructor(props) {
		super(props);
		this.state = {
			busy: false,
			password: '',
			username: ''
		};
		this.nextPathname = this.nextPathname.bind(this);
	}

	componentDidMount() {
		if (AuthService.loggedIn) {
			this.context.router.replace('/');
			return;
		}
	}

	nextPathname() {
		const { location } = this.props;
		if (location.state && location.state.nextPathname) {
			return location.state.nextPathname;
		}
		return '/';
	}

	submitForm(event) {

		event.preventDefault();

		if (!this.state.busy && this.state.username && this.state.password) {
			this.setState({busy: true});
			AuthService.login(this.state.username, this.state.password).then(success => {
				this.setState({busy: false});
				if (success) {
					this.context.router.replace(this.nextPathname());
				} else {
					// TODO: failed login message
				}
			});
		}

	}

	updateForm(event) {
		if (event.target.id === 'username') {
			this.setState({username: event.target.value});
		} else if (event.target.id === 'password') {
			this.setState({password: event.target.value});
		}
	}

	render() {

		return (
			<Card style={style.card}>
				<form onSubmit={(event) => this.submitForm(event)}>

					<CardTitle title="Welcome to Rakewire" subtitle="Login" style={style.cardheader}/>

					<CardText>

						<div style={style.field}>
							<TextField id="username" floatingLabelText="Username"
								value={this.state.username} style={style.input}
								onChange={(event) => this.updateForm(event)} />
						</div>

						<div style={style.field}>
							<TextField id="password" type="password" floatingLabelText="Password"
								value={this.state.password} style={style.input}
								onChange={(event) => this.updateForm(event)} />
						</div>

					</CardText>

					<CardActions>
						<RaisedButton disabled={this.state.busy} label="Login"
							onTouchTap={(event) => this.submitForm(event)} primary={true} type="submit" />
					</CardActions>

				</form>
			</Card>
		);

	}

}

export default withRouter(Login);
