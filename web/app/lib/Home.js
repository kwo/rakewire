import React, { PropTypes } from 'react';

class Home extends React.Component {

	static displayName = 'home';

	static contextTypes = {
		muiTheme : PropTypes.object.isRequired
	};

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		return (
			<div>
				<p>Welcome to Rakewire.</p>
			</div>
		);

	}

}

export default Home;
