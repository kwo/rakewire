/* eslint no-unused-vars: 0 */
import React from 'react';
import {DefaultRoute, Redirect, Route} from 'react-router';
import App from './App';
import About from './About';
import Home from './Home';

// -------------------- routes --------------------

const routes = (
	<Route handler={App} name="app" path="/" >
		<Route handler={Home} name="home" path="/" />
		<Route handler={About} name="about" path="/about" />
		<DefaultRoute handler={Home} />
		<Redirect from="*" to="home" />
	</Route>
);

export default routes;
