import React, {PropTypes} from 'react';
import {Link} from 'react-router';

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		children: PropTypes.node
	};

	constructor(props, context) {
		super(props, context);
	}

	render() {
		return (
			<div>
				<div>
					<Link activeClassName="active" to="/">Home</Link> <Link activeClassName="active" to="about">About</Link>
				</div>
				{this.props.children}
			</div>
		);
	}

}

export default App;
