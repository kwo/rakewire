import React, { PropTypes } from 'react';

class About extends React.Component {

	static displayName = 'about';

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
				<p>About Rakewire</p>
			</div>
		);

	}

}

export default About;
