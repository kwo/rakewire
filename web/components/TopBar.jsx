import React, {PropTypes} from 'react';
import AppBar from 'material-ui/AppBar';
import SideBar from './SideBar';

class TopBar extends React.Component {

	static displayName = 'topbar';

	static contextTypes = {
		router: PropTypes.object.isRequired
	}

	static propTypes = {
		sidebar: PropTypes.bool,
		title: PropTypes.string.isRequired
	};

	static defaultProps = {
		sidebar: false
	}

	constructor(props, ctx) {
		super(props, ctx);
		this.state = {
			sidebar: this.props.sidebar
		};
		this.navigate = this.navigate.bind(this);
		this.toggleSideBar = this.toggleSideBar.bind(this);
	}

	navigate(event, path) {
		this.toggleSideBar(event);
		if (path) this.context.router.push(path);
	}

	toggleSideBar(/*event*/) {
		const state = this.state;
		state.sidebar = !state.sidebar;
		this.setState(state);
	}

	render() {
		return (
			<div>
				<AppBar
					onLeftIconButtonTouchTap={this.toggleSideBar}
					title={this.props.title} />
				<SideBar navigate={this.navigate} opened={this.state.sidebar} title={this.props.title} />
			</div>
		);
	}

}

export default TopBar;
