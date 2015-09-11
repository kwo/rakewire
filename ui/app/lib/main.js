import 'site.css!';
import 'fetch';

import React from 'react';
import Router, {DefaultRoute, Redirect, Route} from 'react-router';
import { Styles } from 'material-ui';
import injectTapEventPlugin from 'npm:react-tap-event-plugin@0.1.7/src/injectTapEventPlugin';

const ThemeManager = new Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.LIGHT);
injectTapEventPlugin();

import About from './About';
import App from './App';
import Home from './Home';
import Feed from './Feed';
import Feeds from './Feeds';

const routes = (
	<Route handler={App} name="app" path="/" >
		<Route handler={Home} name="home" path="/" />
		<Route handler={About} name="about" path="/about" />
		<Route handler={Feeds} name="feeds" path="/feeds" />
		<Route handler={Feed} name="feed" path="/feeds/:id" />
		<DefaultRoute handler={Home} />
		<Redirect from="*" to="home" />
	</Route>
);

Router.run(routes, Router.HistoryLocation, function(Handler) {
	React.render(<Handler />, document.getElementById('app'));
});
