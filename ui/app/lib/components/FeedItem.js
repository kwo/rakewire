import React, { PropTypes } from 'react';
import { TableHeaderColumn, TableRow, TableRowColumn } from 'material-ui';
import moment from 'moment';

// DOING click to detail view

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
		this.onRowClick = this.onRowClick.bind(this);
	}

	onRowClick() {
		console.log(this.props.feed.id);
	}

	render() {

		const style = {
			td: {
				cursor: 'pointer',
				height: 24,
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			},
			tdFeed: {
				cursor: 'pointer',
				height: 24,
				textAlign: 'left',
				width: '50%'
			},
			th: {
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			},
			thFeed: {
				textAlign: 'left',
				width: '50%'
			},
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
					<TableHeaderColumn style={style.thFeed}>Feed</TableHeaderColumn>
				</TableRow>
			);
		}

		const formatDate = function(dt) {
			return moment(dt).format('dd HH:mm');
		};

		const status = feed.last.result;
		const isUpdated = feed.last.updated ? 'Yes' : '';
		const updateCheck = feed.last.updateCheck;
		const title = feed.title || feed.last200.feed.title || feed.url;

		return (
			<TableRow
				hoverable={true}
				onRowClick={this.onRowClick}>
				<TableRowColumn style={style.td}>{formatDate(feed.nextFetch)}</TableRowColumn>
				<TableRowColumn style={style.td}>{formatDate(feed.last.startTime)}</TableRowColumn>
				<TableRowColumn style={style.td}>{status}</TableRowColumn>
				<TableRowColumn style={style.td}>{feed.last.http.statusCode}</TableRowColumn>
				<TableRowColumn style={style.td}>{isUpdated}</TableRowColumn>
				<TableRowColumn style={style.td}>{updateCheck}</TableRowColumn>
				<TableRowColumn style={style.tdFeed}>{title}</TableRowColumn>
			</TableRow>
		);

	}

}

export default FeedItem;
