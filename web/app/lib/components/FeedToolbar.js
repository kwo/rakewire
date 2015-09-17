import React, { PropTypes } from 'react';
import { IconButton, RefreshIndicator, Toolbar, ToolbarGroup, ToolbarTitle } from 'material-ui';
import moment from 'moment';

class FeedToolbar extends React.Component {

	static displayName = 'feedtoolbar';

	static propTypes = {
		lastRefresh: PropTypes.object,
		onRefreshClick: PropTypes.func.isRequired
	};

	static contextTypes = {
		muiTheme : PropTypes.object
	};

	static defaultProps = {
		lastRefresh: null,
		onRefreshClick: function() {
			console.log('FeedToolbar: noop refresh handler');
		}
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		const formatDate = function(dt) {
			if (!dt || dt.valueOf() === 0) return '';
			return moment(dt).format('dd HH:mm:ss');
		};

		const refreshStatus = this.props.lastRefresh ? 'ready' : 'loading';

		return (

			<Toolbar>
				<ToolbarGroup float="left" key={0}>
					<ToolbarTitle text="Feed List" />
				</ToolbarGroup>
				<ToolbarGroup float="right" key={1}>
					<ToolbarTitle text={formatDate(this.props.lastRefresh)} />
					<IconButton onTouchTap={this.props.onRefreshClick} tooltip="Refresh" >
						<RefreshIndicator
							left={10}
							percentage={100}
							size={40}
							status={refreshStatus}
							top={5} />
					</IconButton>
				</ToolbarGroup>
			</Toolbar>

		);

	} // render

}

export default FeedToolbar;
