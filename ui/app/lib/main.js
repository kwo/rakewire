import React from 'react';
import Router, {DefaultRoute, Link, Redirect, Route, RouteHandler} from 'react-router';
import { Styles } from 'material-ui';
import injectTapEventPlugin from 'npm:react-tap-event-plugin@0.1.7/src/injectTapEventPlugin'; // WTF?

import 'site.css!';
const ThemeManager = new Styles.ThemeManager();
injectTapEventPlugin();
ThemeManager.setTheme(ThemeManager.types.LIGHT);

import About from './AboutComponent';
import Home from './HomeComponent';

// -------------------- app component --------------------
class App extends React.Component {

	constructor(props) {
		super(props);
		this.state = {};
	}

	getChildContext () {
		return {
			muiTheme : ThemeManager.getCurrentTheme()
		};
	}

	render() {
		return (
			<div className="container">

				<div className="page-header">
					<Link to="home">Rakewire</Link> | <Link to="about">About</Link>
				</div>

				<RouteHandler />

				<footer className="footer">
					<div className="container">
						<p className="small text-muted">Copyright © 2015 <a href="https://ostendorf.com/">Karl Ostendorf</a></p>
					</div>
				</footer>

			</div>
		);
	}

}
App.displayName = 'app';
// only necessary by App (outermost parent component)
App.childContextTypes = {
	muiTheme : React.PropTypes.object
};


// -------------------- routes --------------------

const routes = (
	<Route handler={App} name="app" path="/" >
		<Route handler={Home} name="home" path="/" />
		<Route handler={About} name="about" path="/about" />
		<DefaultRoute handler={Home} />
		<Redirect from="*" to="home" />
	</Route>
);

Router.run(routes, Router.HistoryLocation, function(Handler) {
	React.render(<Handler />, document.getElementById('app'));
});
