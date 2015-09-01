import React, { PropTypes } from 'react';

class Home extends React.Component {

	static displayName = 'home';

	// static propTypes = {
	// 	max: PropTypes.number
	// };

	static contextTypes = {
		router: PropTypes.func,
		muiTheme : PropTypes.object.isRequired
	};

	static childContextTypes = {
		muiTheme : PropTypes.object
	};

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	// getDefaultProps() {
	// 	return {
	// 		max: 100
	// 	};
	// }

	getChildContext () {
		return {
			muiTheme : this.context.muiTheme
		};
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
