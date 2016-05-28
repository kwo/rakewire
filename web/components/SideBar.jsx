import React, {PropTypes} from 'react';
import AuthService from '../services/Auth';

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import {List, ListItem, MakeSelectable} from 'material-ui/List';

import DashboardIcon from 'material-ui/svg-icons/action/dashboard';
import ExitIcon from 'material-ui/svg-icons/action/exit-to-app';
import FeedIcon from 'material-ui/svg-icons/communication/rss-feed';
import HomeIcon from 'material-ui/svg-icons/action/home';
import IconButton from 'material-ui/IconButton';

class SideBar extends React.Component {

	static displayName = 'sidebar';

	static propTypes = {
		navigate: PropTypes.func.isRequired,
		opened: PropTypes.bool,
		title: PropTypes.string.isRequired
	};

	static defaultProps = {
		opened: false
	}

	constructor(props, ctx) {
		super(props, ctx);
		this.state = {};
	}

	render() {

		const SelectableList = MakeSelectable(List);

		const iconApp=(
			<IconButton disabled={true}>
				<FeedIcon/>
			</IconButton>
		);

		return (
			<Drawer
				docked={false}
				onRequestChange={(event) => this.props.navigate(event)}
				open={this.props.opened} >

				<AppBar iconElementLeft={iconApp} title={this.props.title} />

				<SelectableList onChange={this.props.navigate}>
					<ListItem primaryText="Home" value="home" leftIcon={<HomeIcon />} />
					<ListItem primaryText="Dashboard" value="dashboard" leftIcon={<DashboardIcon />} />
					{AuthService.loggedIn && (
					<ListItem primaryText="Logout" value="logout" leftIcon={<ExitIcon />} />
					)}
				</SelectableList>

			</Drawer>
		);
	}

}

export default SideBar;
