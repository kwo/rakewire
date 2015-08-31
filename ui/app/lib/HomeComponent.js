import React from 'react';
import mui from 'material-ui';

const FloatingActionButton = mui.FloatingActionButton;
const Colors = mui.Styles.Colors;
const FontIcon = mui.FontIcon;


class Home extends React.Component {

	constructor(props) {
		super(props);
		this.state = {};
	}

	render() {

		return (
			<div>

				<p>Welcome to Rakewire.</p>
				<FloatingActionButton mini={true} secondary={true} >
					<FontIcon className="material-icons" hoverColor={Colors.red500}>star</FontIcon>
				</FloatingActionButton>

			</div>
		);

	}

}

Home.displayName = 'home';

export default Home;
