import React, {PropTypes} from 'react';
import { withRouter } from 'react-router';
import AuthService from '../services/Auth';

class Logout extends React.Component {

	static displayName = 'logout';

	static propTypes = {
		location: PropTypes.object,
		router: PropTypes.object.isRequired
	};

	constructor(props) {
		super(props);
		this.state = {};
	}

	componentDidMount() {

		if (!AuthService.loggedIn) {
			this.props.router.replace('/');
			return;
		}

		AuthService.logout().then(() => {
			const { location } = this.props;
			if (location.state && location.state.nextPathname) {
				this.props.router.replace(location.state.nextPathname);
			} else {
				this.props.router.replace('/');
			}
		});

	}

	render() {
		return (<div></div>);
	}

}

export default withRouter(Logout);
