import React from 'react';
import ReactDOM from 'react-dom';
import Router, {IndexRoute, Redirect, Route} from 'react-router';
import createBrowserHistory from 'history/lib/createBrowserHistory';

import About from './About';
import App from './App';
import Home from './Home';
import Feed from './Feed';
import Feeds from './Feeds';

const routes = (
	<Router history={createBrowserHistory()}>
		<Route component={App} path="/">
			<IndexRoute component={Home} />
			<Route component={About} path="about" />
			<Route component={Feeds} path="feeds" />
			<Route component={Feed}  path="feeds/:id" />
			<Redirect from="*" to="/" />
		</Route>
	</Router>
);

ReactDOM.render(routes, document.getElementById('app'));
