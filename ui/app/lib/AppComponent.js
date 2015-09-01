import React from 'react';
import { Link, RouteHandler } from 'react-router';
import { AppBar, IconButton, Styles, Tab, Tabs } from 'material-ui';
const ThemeManager = new Styles.ThemeManager();

class App extends React.Component {

	static displayName = 'app';

	static contextTypes = {
		router: React.PropTypes.func
	};

	// only necessary by App (outermost parent component)
	static childContextTypes = {
		muiTheme : React.PropTypes.object
	};

	constructor(props, context) {

		super(props, context);

		const currentRoutes = this.context.router.getCurrentRoutes();
		const activeRouteName = currentRoutes[currentRoutes.length - 1].name;
		this.state = {
			tab: activeRouteName
		};

		this.onLogoClick = this.onLogoClick.bind(this);
		this.onChangeTabs = this.onChangeTabs.bind(this);
		this.navigateTo = this.navigateTo.bind(this);

	}

	getChildContext () {
		return {
			muiTheme : ThemeManager.getCurrentTheme()
		};
	}

	navigateTo(name) {
		const state = this.state;
		state.tab = name;
		this.setState(state);
		this.context.router.transitionTo(state.tab);
	}

	onChangeTabs(name /*, event, tab*/) {
		this.navigateTo(name);
	}

	onLogoClick(/*event*/) {
		this.navigateTo('home');
	}

	render() {

		const styles = {
			appBar: {
				flexWrap: 'wrap',
			},
			tabs: {
				width: '25%',
			}
		};


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
					iconElementLeft={logoButton}
					style={styles.appBar}
					title="Rakewire">
					<Tabs
						onChange={this.onChangeTabs}
						style={styles.tabs}
						value={this.state.tab}>
						<Tab label="Home" value="home" />
						<Tab label="About" value="about" />
					</Tabs>
				</AppBar>

				<RouteHandler />

			</div>
		);

	}

}

export default App;
