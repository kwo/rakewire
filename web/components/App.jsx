import React, {PropTypes} from 'react';
import {Link} from 'react-router';
import AuthService from '../services/Auth';

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		children: PropTypes.node
	};

	constructor(props) {
		super(props);
	}

	render() {
		return (
			<div>
				<div>
					<Link activeClassName="active" to="/">Home</Link><span> </span>
					<Link activeClassName="active" to="about">About</Link><span> </span>
					{AuthService.loggedIn && (<Link activeClassName="active" to="logout">Logout</Link>)}
				</div>
				{this.props.children}
			</div>
		);
	}

}

export default App;
