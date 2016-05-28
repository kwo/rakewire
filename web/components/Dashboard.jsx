import React, {PropTypes} from 'react';

class Dashboard extends React.Component {

	static displayName = 'dashboard';

	static propTypes = {
		children: PropTypes.node
	};

	static defaultProps = {
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {
		return (
			<p>Dashboard</p>
		);
	}

}

export default Dashboard;
