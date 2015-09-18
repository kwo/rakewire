import React from 'react';
import {Link} from 'react-router';
import {CollapsibleNav, Nav, Navbar} from 'react-bootstrap';
import NavLink from './NavLink';

class AppHeader extends React.Component {

	static displayName = 'appheader';

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		return (
			<header>
				<Navbar brand={<Link to="home">Rakewire</Link>} fluid={true} toggleNavKey={0}>
					<CollapsibleNav eventKey={0} fluid={true} >
						<Nav navbar>
							<NavLink to="feeds">Feeds</NavLink>
						</Nav>
						<Nav navbar right>
							<NavLink to="about">About</NavLink>
						</Nav>
					</CollapsibleNav>
				</Navbar>
			</header>
		);

	} // render

}

export default AppHeader;
