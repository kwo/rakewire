import React from 'react';

class Home extends React.Component {

	constructor(props) {
		super(props);
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
