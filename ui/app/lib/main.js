import 'site.css!';
import React from 'react';
import Router, {DefaultRoute, Redirect, Route} from 'react-router';
import { Styles } from 'material-ui';
import injectTapEventPlugin from 'npm:react-tap-event-plugin@0.1.7/src/injectTapEventPlugin'; // WTF?

const ThemeManager = new Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.LIGHT);
injectTapEventPlugin();

import App from './AppComponent';
import About from './AboutComponent';
import Home from './HomeComponent';

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
