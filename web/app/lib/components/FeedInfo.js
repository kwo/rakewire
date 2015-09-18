import React, { PropTypes } from 'react';
//import moment from 'moment';

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

		// const formatDate = function(dt) {
		// 	return moment(dt).format('dd HH:mm');
		// };

		const title = feed.title || feed.last200.feed.title || feed.url;

		return (
			<thead>
				<tr key={0}>
					<td><strong>Title</strong></td>
					<td>{title}</td>
				</tr>

				<tr key={1}>
					<td><strong>URL</strong></td>
					<td>{feed.url}</td>
				</tr>
			</thead>
		);

	}

}

export default FeedInfo;
