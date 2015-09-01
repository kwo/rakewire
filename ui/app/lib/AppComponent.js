import React from 'react';
import {Link, RouteHandler} from 'react-router';
import { Styles } from 'material-ui';
const ThemeManager = new Styles.ThemeManager();

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

export default App;
