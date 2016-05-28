// root directory files
import 'file?name=[name].[ext]!./robots.txt';
import 'file?name=[name].[ext]!./favicon.ico';

// css
import './node_modules/roboto-fontface/css/roboto-fontface-regular.scss';

// import global modules
import 'whatwg-fetch';
import injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();

import React from 'react';
import ReactDOM from 'react-dom';
import {browserHistory, IndexRoute, Redirect, Route, Router} from 'react-router';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

import theme from './theme';

// services
import AuthService from './services/Auth';

// components
import About from './components/About';
import App from './components/App';
import Dashboard from './components/Dashboard';
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
				<Route component={Dashboard} path="dashboard" />
			</Route>
			<Redirect from="*" to="/" />
		</Route>
	</Router>
);

const app = (
	<MuiThemeProvider muiTheme={getMuiTheme(theme)}>
		{routes}
	</MuiThemeProvider>
);

ReactDOM.render(app, document.getElementById('app'));
