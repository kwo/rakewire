import React, { PropTypes } from 'react';
import { TableHeaderColumn, TableRow, TableRowColumn } from 'material-ui';
import moment from 'moment';

// #DOING:70 click to detail view

class FeedItem extends React.Component {

	static displayName = 'feeditem';

	static propTypes = {
		feed: PropTypes.object,
		hoverCol: PropTypes.number,
		hovered: PropTypes.bool
	};

	static contextTypes = {
		muiTheme : PropTypes.object
	};

	static defaultProps = {
		feed: null,
		hoverCol: 0,
		hovered: false
	}

	constructor(props, context) {
		super(props, context);
		this.state = {
			hoverCol: this.props.hoverCol,
			hovered: this.props.hovered
		};
		this.onCellClick = this.onCellClick.bind(this);
		this.onCellHover = this.onCellHover.bind(this);
		this.onCellHoverExit = this.onCellHoverExit.bind(this);
	}

	onCellClick(/*e, id*/) {
		console.log(this.props.feed);
	}

	onCellHover(e, row, col) {
		const state = this.state;
		state.hovered = true;
		state.hoverCol = col;
		this.setState(state);
	}

	onCellHoverExit(e, row, col) {
		const state = this.state;
		state.hovered = false;
		state.hoverCol = col;
		this.setState(state);
	}

	render() {

		const style = {
			td: {
				height: 24,
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			},
			tdFeed: {
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
					<TableHeaderColumn style={style.th} tooltip="time of the next scheduled fetch">Next</TableHeaderColumn>
					<TableHeaderColumn style={style.th} tooltip="the time of the last attempted fetch">Last</TableHeaderColumn>
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

		const formatMessage = function(f) {
			return f.last.resultMessage ? ': ' + f.last.resultMessage : '';
		}

		const formatStatus = function(value) {
			switch (value) {
			case 'OK':
				return 'OK';
			case 'MV':
				return 'MV = redirected';
			case 'EC':
				return 'EC = error client' + formatMessage(feed);
			case 'ES':
				return 'ES = error server' + formatMessage(feed);
			case 'FP':
				return 'FP = cannot parse feed' + formatMessage(feed);
			case 'FT':
				return 'FT = cannot parse feed time' + formatMessage(feed);
			}
			return value + ' = unknown';
		};

		const formatUpdateCheck = function(value) {
			switch (value) {
			case 'LU':
				return 'LU = Last Updated';
			case 'NM':
				return 'NM = Not Modifed';
			}
			return value + ' = Unknown';
		};

		let status = feed.last.result;
		let isUpdated = feed.last.updated ? 'Yes' : '';
		let updateCheck = feed.last.updateCheck;
		let title = feed.title || feed.last200.feed.title || feed.url;

		if (this.state.hovered) {
			switch (this.state.hoverCol) {
			case 3:
				title = formatStatus(feed.last.result);
				break;
			case 6:
				title = formatUpdateCheck(feed.last.updateCheck);
				break;
			case 7:
				title = feed.url;
				break;
			}
		}

		return (
			<TableRow
				hoverable={true}
				key={feed.id}
				onCellHover={this.onCellHover}
				onCellHoverExit={this.onCellHoverExit}
				onClick={this.onCellClick}>
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
