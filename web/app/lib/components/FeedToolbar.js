import React, { PropTypes } from 'react';
import moment from 'moment';

class FeedToolbar extends React.Component {

	static displayName = 'feedtoolbar';

	static propTypes = {
		lastRefresh: PropTypes.object,
		onRefreshClick: PropTypes.func.isRequired
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

			<div>
				<button onClick={this.props.onRefreshClick} tooltip="Refresh" >
				</button>
			</div>

		);

	} // render

}

export default FeedToolbar;
