import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import moment from 'moment';

// TODO: eliminate last and last200 so that everything is in feed

class FeedEntry extends React.Component {

	static displayName = 'feedentry';

	static propTypes = {
		feed: PropTypes.object,
		lastRefresh: PropTypes.object,
		onRefreshClick: PropTypes.func
	};

	static defaultProps = {
		feed: null,
		lastRefresh: null,
		onRefreshClick: function() {}
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		const displayRefresh = function(dt) {
			if (!dt) return (<i className="material-icons">loop</i>);
			return formatTime(dt);
		};

		// const formatDate = function(dt) {
		// 	return moment(dt).format('YYYY-MM-DD');
		// };

		const formatTime = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const feed = this.props.feed;

		if (!feed) {
			return (
				<tr key={0}>
					<th>Status</th>
					<th>Next</th>
					<th>
						Feed
						<span className="pull-right"
									onClick={this.props.onRefreshClick}
									style={{cursor: 'pointer'}}
									tooltip="Refresh">
							{displayRefresh(this.props.lastRefresh)}
						</span>
					</th>
				</tr>
			);
		}

		const status = feed.status;
		const title = feed.title || feed.url;
		const feedLinkURL = `/feeds/${feed.id}`;

		return (
			<tr>
				<td>{status}</td>
				<td>{formatTime(feed.nextFetch)}</td>
				<td><Link to={feedLinkURL}>{title}</Link></td>
			</tr>
		);

	}

}

export default FeedEntry;
