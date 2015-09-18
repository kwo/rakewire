import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import moment from 'moment';

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

		const formatDate = function(dt) {
			return moment(dt).format('YYYY-MM-DD');
		};

		const formatTime = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const feed = this.props.feed;

		if (!feed) {
			return (
				<tr key={0}>
					<th>Next</th>
					<th>Last</th>
					<th>Status</th>
					<th>Code</th>
					<th>Updated</th>
					<th>Check</th>
					<th>Last Updated</th>
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

		const status = feed.last.result;
		const isUpdated = feed.last.updated ? 'Yes' : '';
		const updateCheck = feed.last.updateCheck;
		const title = feed.title || feed.last200.feed.title || feed.url;
		const feedLinkURL = '/feeds/' + feed.id;

		return (
			<tr>
				<td>{formatTime(feed.nextFetch)}</td>
				<td>{formatTime(feed.last.startTime)}</td>
				<td>{status}</td>
				<td>{feed.last.http.statusCode}</td>
				<td>{isUpdated}</td>
				<td>{updateCheck}</td>
				<td>{formatDate(feed.lastUpdated)}</td>
				<td><Link to={feedLinkURL}>{title}</Link></td>
			</tr>
		);

	}

}

export default FeedEntry;
