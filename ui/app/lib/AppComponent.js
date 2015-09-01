import React from 'react';
import { Link, RouteHandler } from 'react-router';
import { AppBar, IconButton, Styles } from 'material-ui';
const ThemeManager = new Styles.ThemeManager();

class App extends React.Component {

		this.state = {};
	constructor(opts) {
		super(opts);
	}

	getChildContext () {
		return {
			muiTheme : ThemeManager.getCurrentTheme()
		};
	}

	onLogoClick(event) {
		console.log(event)
	}

	render() {

		const logoButton = (
			<IconButton
				containerElement={<Link to="/home" />}
				iconClassName="material-icons"
				linkButton={true}
				onTouchTap={this.onLogoClick}>
				star
			</IconButton>
		);

		return (
			<div>
				<AppBar
					title="Rakewire"
					iconElementLeft={logoButton}
					/>
				<RouteHandler />
			</div>
		);

	}

}

App.displayName = 'app';

// only necessary by App (outermost parent component)
App.childContextTypes = {
	muiTheme : React.PropTypes.object
};

export default App;
