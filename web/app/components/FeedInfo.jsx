import React, { PropTypes } from 'react';
import moment from 'moment';

class FeedInfo extends React.Component {

	static displayName = 'feedinfo';

	static propTypes = {
		feed: PropTypes.object
	};

	static defaultProps = {
		feed: null
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		const feed = this.props.feed;

		if (!feed) {
			return (
				<div>no feed</div>
			);
		}

		const formatDate = function(dt) {
			return moment(dt).format('YYYY-MM-DD HH:mm');
		};

		const title = feed.title || feed.url;

		return (
			<thead>

				<tr key={0}>
					<td><strong>ID</strong></td>
					<td>{feed.id}</td>
				</tr>

				<tr key={1}>
					<td><strong>Title</strong></td>
					<td>{title}</td>
				</tr>

				<tr key={2}>
					<td><strong>URL</strong></td>
					<td><a href={feed.url}>{feed.url}</a></td>
				</tr>

				<tr key={3}>
					<td><strong>Last Updated</strong></td>
					<td>{formatDate(feed.lastUpdated)}</td>
				</tr>

				<tr key={4}>
					<td><strong>Next Fetch</strong></td>
					<td>{formatDate(feed.nextFetch)}</td>
				</tr>

				<tr key={5}>
					<td><strong>Notes</strong></td>
					<td>{feed.notes}</td>
				</tr>

			</thead>
		);

	}

}

export default FeedInfo;
