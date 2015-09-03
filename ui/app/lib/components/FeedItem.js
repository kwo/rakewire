import React, { PropTypes } from 'react';
import moment from 'moment';

// #DOING:20 add style to row element (gray text, more spacing)
// #DOING:30 tooltips for abbr
// #DOING:70 click to detail view

class FeedItem extends React.Component {

	static displayName = 'feeditem';

	static propTypes = {
		feed: PropTypes.object
	};

	static contextTypes = {
		muiTheme : PropTypes.object
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
				<tr key={0}>
					<th>Next</th>
					<th>Last</th>
					<th>Status</th>
					<th>Code</th>
					<th>Updated</th>
					<th>Check</th>
					<th>Feed</th>
				</tr>
			);
		}

		const formatDate = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const title = feed.title || feed.last200.title || feed.url;
		const isUpdated = feed.last.updated ? 'Yes' : '';

		return (
			<tr key={feed.id}>
				<td>{formatDate(feed.nextFetch)}</td>
				<td>{formatDate(feed.last.startTime)}</td>
				<td>{feed.last.result}</td>
				<td>{feed.last.http.statusCode}</td>
				<td>{isUpdated}</td>
				<td>{feed.last.updateCheck}</td>
				<td>{title}</td>
			</tr>
		);

	}

}

export default FeedItem;
