import React, { PropTypes } from 'react';

class About extends React.Component {

	static displayName = 'about';

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
				<p>About Rakewire</p>
			</div>
		);

	}

}

export default About;
