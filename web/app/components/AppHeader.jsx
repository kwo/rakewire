import React from 'react';
import {IndexLink} from 'react-router';
import {CollapsibleNav, Nav, Navbar, NavBrand} from 'react-bootstrap';
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
				<Navbar fluid={true}
								inverse={true}
								toggleNavKey={0}>
					<NavBrand>{<IndexLink to="/">Rakewire</IndexLink>}</NavBrand>
					<CollapsibleNav eventKey={0} fluid={true} >
						<Nav navbar>
							<NavLink to="/feeds">Feeds</NavLink>
						</Nav>
						<Nav navbar right>
							<NavLink to="/about">About</NavLink>
						</Nav>
					</CollapsibleNav>
				</Navbar>
			</header>
		);

	} // render

}

export default AppHeader;
