import React from 'react';

class About extends React.Component {

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

About.displayName = 'about';

export default About;
