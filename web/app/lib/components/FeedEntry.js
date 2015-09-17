import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { TableHeaderColumn, TableRow, TableRowColumn } from 'material-ui';
import moment from 'moment';

class FeedEntry extends React.Component {

	static displayName = 'feedentry';

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
				<TableRow key={0}>
					<TableHeaderColumn style={style.th}>Next</TableHeaderColumn>
					<TableHeaderColumn style={style.th}>Last</TableHeaderColumn>
					<TableHeaderColumn style={style.th}>Status</TableHeaderColumn>
					<TableHeaderColumn style={style.th}>Code</TableHeaderColumn>
					<TableHeaderColumn style={style.th}>Updated</TableHeaderColumn>
					<TableHeaderColumn style={style.th}>Check</TableHeaderColumn>
					<TableHeaderColumn style={style.thLU}>Last Updated</TableHeaderColumn>
					<TableHeaderColumn style={style.thFeed}>Feed</TableHeaderColumn>
				</TableRow>
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
			<TableRow
				hoverable={true}
				onRowClick={this.onRowClick}>
				<TableRowColumn style={style.td}>{formatTime(feed.nextFetch)}</TableRowColumn>
				<TableRowColumn style={style.td}>{formatTime(feed.last.startTime)}</TableRowColumn>
				<TableRowColumn style={style.td}>{status}</TableRowColumn>
				<TableRowColumn style={style.td}>{feed.last.http.statusCode}</TableRowColumn>
				<TableRowColumn style={style.td}>{isUpdated}</TableRowColumn>
				<TableRowColumn style={style.td}>{updateCheck}</TableRowColumn>
				<TableRowColumn style={style.tdLU}>{formatDate(feed.lastUpdated)}</TableRowColumn>
				<TableRowColumn style={style.tdFeed}><Link to={feedLinkURL}>{title}</Link></TableRowColumn>
			</TableRow>
		);

	}

}

export default FeedEntry;
