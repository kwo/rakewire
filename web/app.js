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

// services
import AuthService from './services/Auth';

// components
import About from './components/About';
import App from './components/App';
import Home from './components/Home';
import Login from './components/Login';
import Logout from './components/Logout';

function loginIfUnauthenticated(nextState, replace) {
	if (!AuthService.loggedIn) {
		replace({
			pathname: 'login',
			state: { nextPathname: nextState.location.pathname }
		});
	}
}

const routes = (
	<Router history={browserHistory}>
		<Route component={App} >
			<Route component={Login}  path="login" />
			<Route component={Logout} path="logout" />
			<Route component={About}  path="about" />
			<Route onEnter={loginIfUnauthenticated} path="/" >
				<IndexRoute component={Home} />
			</Route>
			<Redirect from="*" to="/" />
		</Route>
	</Router>
);

ReactDOM.render(routes, document.getElementById('app'));
