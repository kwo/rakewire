import React, {PropTypes} from 'react';
import { RouteHandler } from 'react-router';
import { AppBar, IconButton, Styles, Tab, Tabs } from 'material-ui';
const ThemeManager = new Styles.ThemeManager();

import TitleComponent from './components/Title';

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		title: PropTypes.string
	};

	static contextTypes = {
		router: PropTypes.func
	};

	static childContextTypes = {
		muiTheme : PropTypes.object
	};

	static defaultProps = {
		title: 'Rakewire'
	}

	constructor(props, context) {

		super(props, context);

		this.state = {
			tab: this.context.router.getLocation().getCurrentPath()
		};

		this.navigateTo = this.navigateTo.bind(this);
		this.onChangeTabs = this.onChangeTabs.bind(this);
		this.onLogoClick = this.onLogoClick.bind(this);
		this.onRouteChange = this.onRouteChange.bind(this);

		this.context.router.getLocation().addChangeListener(this.onRouteChange);

	}

	getChildContext () {
		return {
			muiTheme: ThemeManager.getCurrentTheme()
		};
	}

	navigateTo(path) {
		const currentPath = this.context.router.getLocation().getCurrentPath();
		if (currentPath !== path) {
			this.context.router.transitionTo(path);
		}
	}

	onChangeTabs(path /*, event, tab*/) {
		this.navigateTo(path);
	}

	onLogoClick(/*event*/) {
		this.navigateTo('/');
	}

	onRouteChange(event) {
		// types: pop, push, replace
		// console.log(event);
		const state = this.state;
		state.tab = event.path;
		this.setState(state);
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
					title={<TitleComponent onClick={this.onLogoClick} title={this.props.title} />}>
					<Tabs
						onChange={this.onChangeTabs}
						style={styles.tabs}
						value={this.state.tab}>
						<Tab label="Home" value="/" />
						<Tab label="About" value="/about" />
					</Tabs>
				</AppBar>

				<RouteHandler />

			</div>
		);

	}

}

export default App;
