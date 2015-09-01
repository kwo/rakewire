import React from 'react';

class Home extends React.Component {

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

Home.displayName = 'home';

export default Home;
