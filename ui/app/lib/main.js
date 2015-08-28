import React from 'react';
import Router from 'react-router';

import About from './AboutComponent';
import Home from './HomeComponent';

const Route = Router.Route;
const DefaultRoute = Router.DefaultRoute;
const Link = Router.Link;
const Redirect = Router.Redirect;
const RouteHandler = Router.RouteHandler;

// -------------------- app component --------------------

class App extends React.Component {

	constructor(props) {
		super(props);
		this.state = {};
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
