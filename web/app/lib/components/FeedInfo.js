import React, { PropTypes } from 'react';
import moment from 'moment';

class FeedInfo extends React.Component {

	static displayName = 'feedinfo';

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
				<div>no feed</div>
			);
		}

		const formatDate = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const style = {
			td: {
				height: 24,
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			},
			th: {
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			}
		};

		const title = feed.title || feed.last200.feed.title || feed.url;

		return (
			<thead>
				<tr key={0}>
					<th style={style.th}>Title</th>
					<td style={style.td}>{title}</td>
				</tr>

				<tr key={1}>
					<th style={style.th}>URL</th>
					<td style={style.td}>{feed.url}</td>
				</tr>
			</thead>
		);

	}

}

export default FeedInfo;
