import '../site.css';
import 'whatwg-fetch';

import React from 'react';
import Router, {DefaultRoute, Redirect, Route} from 'react-router';

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

Router.run(routes, Router.HistoryLocation, Handler => {
	React.render(<Handler />, document.getElementById('app'));
});
