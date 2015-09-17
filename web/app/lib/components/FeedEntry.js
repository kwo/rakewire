import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import moment from 'moment';

class FeedEntry extends React.Component {

	static displayName = 'feedentry';

	static propTypes = {
		feed: PropTypes.object
	};

	static defaultProps = {
		feed: null
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
		this.onRowClick = this.onRowClick.bind(this);
	}

	onRowClick() {
		console.log(this.props.feed.id);
	}

	render() {

		const style = {
			th: {
				paddingLeft: 6,
				paddingRight: 6,
				textAlign: 'left',
				width: '40px'
			},
			thLU: {
				paddingLeft: 6,
				paddingRight: 6,
				textAlign: 'left',
				width: '50px'
			},
			thFeed: {
				paddingLeft: 6,
				paddingRight: 6,
				textAlign: 'left',
				width: '40%'
			},
			td: {
				height: 24,
				paddingLeft: 6,
				paddingRight: 6,
				textAlign: 'left',
				width: '40px'
			},
			tdLU: {
				height: 24,
				paddingLeft: 6,
				paddingRight: 6,
				textAlign: 'left',
				width: '50px'
			},
			tdFeed: {
				height: 24,
				paddingLeft: 6,
				paddingRight: 6,
				textAlign: 'left',
				width: '40%'
			}
		};

		const feed = this.props.feed;

		if (!feed) {
			return (
				<tr key={0}>
					<th style={style.th}>Next</th>
					<th style={style.th}>Last</th>
					<th style={style.th}>Status</th>
					<th style={style.th}>Code</th>
					<th style={style.th}>Updated</th>
					<th style={style.th}>Check</th>
					<th style={style.thLU}>Last Updated</th>
					<th style={style.thFeed}>Feed</th>
				</tr>
			);
		}

		const formatDate = function(dt) {
			return moment(dt).format('YYYY-MM-DD');
		};

		const formatTime = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const status = feed.last.result;
		const isUpdated = feed.last.updated ? 'Yes' : '';
		const updateCheck = feed.last.updateCheck;
		const title = feed.title || feed.last200.feed.title || feed.url;
		const feedLinkURL = '/feeds/' + feed.id;

		return (
			<tr onClick={this.onRowClick}>
				<td style={style.td}>{formatTime(feed.nextFetch)}</td>
				<td style={style.td}>{formatTime(feed.last.startTime)}</td>
				<td style={style.td}>{status}</td>
				<td style={style.td}>{feed.last.http.statusCode}</td>
				<td style={style.td}>{isUpdated}</td>
				<td style={style.td}>{updateCheck}</td>
				<td style={style.tdLU}>{formatDate(feed.lastUpdated)}</td>
				<td style={style.tdFeed}><Link to={feedLinkURL}>{title}</Link></td>
			</tr>
		);

	}

}

export default FeedEntry;
