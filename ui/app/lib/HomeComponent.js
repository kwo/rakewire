import React from 'react';
import mui from 'material-ui';

const RaisedButton = mui.RaisedButton;

class Home extends React.Component {

	constructor(props) {
		super(props);
		this.state = {};
	}

	render() {

		return (
			<div>

				<p>Welcome to Rakewire.</p>
				<RaisedButton label="Hello"/>

			</div>
		);

	}

}

Home.displayName = 'home';

export default Home;
