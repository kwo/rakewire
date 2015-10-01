import React, {PropTypes} from 'react';
import { RouteHandler } from 'react-router';
import AppHeader from '../components/AppHeader';

const Config = {
	rootURL: '/api'
};

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		title: PropTypes.string
	};

	static contextTypes = {
		router: PropTypes.func.isRequired
	};

	static childContextTypes = {
		config: PropTypes.object
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
			config: Config
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

		return (
			<div>
				<AppHeader />
				<div className="container-fluid">
					<RouteHandler />
				</div>
			</div>
		);

	}

}

export default App;
