import React, { PropTypes } from 'react';
import { TableHeaderColumn, TableRow, TableRowColumn } from 'material-ui';
import moment from 'moment';

// #DOING:30 tooltips for abbr
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
		this.onCellHover = this.onCellHover.bind(this);
		this.onCellHoverExit = this.onCellHoverExit.bind(this);
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

		const title = (this.state.hovered && this.state.hoverCol == 7) ? feed.url : feed.title || feed.last200.feed.title || feed.url;
		const isUpdated = feed.last.updated ? 'Yes' : '';

		return (
			<TableRow
				hoverable={true}
				key={feed.id}
				onCellHover={this.onCellHover}
				onCellHoverExit={this.onCellHoverExit}>
				<TableRowColumn style={style.td}>{formatDate(feed.nextFetch)}</TableRowColumn>
				<TableRowColumn style={style.td}>{formatDate(feed.last.startTime)}</TableRowColumn>
				<TableRowColumn style={style.td}>{feed.last.result}</TableRowColumn>
				<TableRowColumn style={style.td}>{feed.last.http.statusCode}</TableRowColumn>
				<TableRowColumn style={style.td}>{isUpdated}</TableRowColumn>
				<TableRowColumn style={style.td}>{feed.last.updateCheck}</TableRowColumn>
				<TableRowColumn style={style.tdFeed}>{title}</TableRowColumn>
			</TableRow>
		);

	}

}

export default FeedItem;
