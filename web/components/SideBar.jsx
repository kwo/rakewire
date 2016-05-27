import React, {PropTypes} from 'react';
import AppBar from 'material-ui/AppBar';
import AuthService from '../services/Auth';
import Drawer from 'material-ui/Drawer';
import FeedIcon from 'material-ui/svg-icons/communication/rss-feed';
import IconButton from 'material-ui/IconButton';
import {List, ListItem, MakeSelectable} from 'material-ui/List';

// TODO: active class names links

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

	// componentDidMount() {
	// }
	//
	// componentWillUnmount() {
	// }

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
					<ListItem primaryText="Home" value="home" />
					<ListItem primaryText="About" value="about" />
					{AuthService.loggedIn && (
						<ListItem primaryText="Logout" value="logout" />
					)}
				</SelectableList>

			</Drawer>
		);
	}

}

export default SideBar;
