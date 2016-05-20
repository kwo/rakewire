import React from 'react';
import {Link} from 'react-router';
import AuthService from '../services/Auth';
import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import IconButton from 'material-ui/IconButton';
import MenuItem from 'material-ui/MenuItem';
import NavigationClose from 'material-ui/svg-icons/navigation/close';

// TODO: active class names links

class TopBar extends React.Component {

	static displayName = 'topbar';

	constructor(props) {
		super(props);
		this.state = {
			sidebar: false // TODO: props
		};
	}

	toggleSidebar(/*event*/) {
		const state = this.state;
		state.sidebar = !state.sidebar;
		this.setState(state);
	}

	render() {
		const iconClose=(
			<IconButton
				onTouchTap={(event) => this.toggleSidebar(event)} >
				<NavigationClose/>
			</IconButton>
		);
		return (
			<div>
				<AppBar
					onLeftIconButtonTouchTap={(event) => this.toggleSidebar(event)}
					title="Rakewire" />
				<Drawer
					docked={false}
					onRequestChange={(event) => this.toggleSidebar(event)}
					open={this.state.sidebar}
					width={200} >
					<AppBar
						iconElementLeft={iconClose}
						onLeftIconButtonTouchTap={(event) => this.toggleSidebar(event)}
						onTitleTouchTap={(event) => this.toggleSidebar(event)}
						title="Rakewire" />

						<MenuItem onTouchTap={(event) => this.toggleSidebar(event)}>
							<Link activeClassName="active" to="/">Home</Link>
						</MenuItem>
						<MenuItem onTouchTap={(event) => this.toggleSidebar(event)}>
							<Link activeClassName="active" to="about">About</Link>
						</MenuItem>
						{AuthService.loggedIn && (
							<MenuItem onTouchTap={(event) => this.toggleSidebar(event)}>
								<Link activeClassName="active" to="logout">Logout</Link>
							</MenuItem>
						)}

				</Drawer>
			</div>
		);
	}

}

export default TopBar;
