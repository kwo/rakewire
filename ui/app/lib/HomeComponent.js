import React from 'react';
import ReactMDL from 'react-mdl';

const Button = ReactMDL.Button;

class Home extends React.Component {

	constructor(props) {
		super(props);
		this.state = {};
	}

	render() {

		return (
			<div>

				<p>Welcome to Rakewire.</p>
				<Button raised={true} colored={true} ripple={false}>Button</Button>

			</div>
		);

	}

}

Home.displayName = 'home';

export default Home;
