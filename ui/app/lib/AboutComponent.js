import React from 'react';

class About extends React.Component {

	static displayName = 'about';

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
