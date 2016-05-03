// root directory files
import 'file?name=[name].[ext]!./robots.txt';

// css
import './node_modules/normalize.css/normalize.css';
import './node_modules/milligram/dist/milligram.css';
import './app.css';

// import global modules
import 'whatwg-fetch';

import React from 'react';
import ReactDOM from 'react-dom';
import {browserHistory, IndexRoute, Redirect, Route, Router} from 'react-router';

import About from './components/About';
import App from './components/App';
import Home from './components/Home';

const routes = (
	<Router history={browserHistory}>
		<Route component={App} path="/">
			<IndexRoute component={Home} />
			<Route component={About} path="about" />
			<Redirect from="*" to="/" />
		</Route>
	</Router>
);

ReactDOM.render(routes, document.getElementById('app'));
