import React from 'react';

class Home extends React.Component {

	static displayName = 'home';

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
