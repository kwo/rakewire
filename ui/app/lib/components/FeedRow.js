import React, { PropTypes } from 'react';
import moment from 'moment';

// #DOING:10 add style to row element (gray text, more spacing)
// #DOING:10 tooltips for abbr
// #DOING:10 click to detail view

class FeedRow extends React.Component {

	static displayName = 'feedrow';

	static propTypes = {
		feed: PropTypes.object.isRequired
	};

	static contextTypes = {
		muiTheme : PropTypes.object
	};

	static defaultProps = {
		feed: {}
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		const formatDate = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const feed = this.props.feed;
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
		)
	}

}

export default FeedRow;
