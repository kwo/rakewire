import React, { PropTypes } from 'react';
import { Table, TableBody, TableHeader } from 'material-ui';
import agent from 'superagent';
import FeedItem from './components/FeedItem';
import FeedToolbar from './components/FeedToolbar';

// #DOING:30 lose state after route change - save to app-wide repository - localstate
// #DOING:40 auto-reload if state too old

class Feeds extends React.Component {

	static displayName = 'feeds';

	static propTypes = {
		feeds: PropTypes.array,
		lastRefresh: PropTypes.object
	};

	static contextTypes = {
		config: PropTypes.object.isRequired,
		muiTheme: PropTypes.object.isRequired
	};

	static defaultProps = {
		feeds: [],
		lastRefresh: new Date(0)
	}

	constructor(props, context) {
		super(props, context);
		this.state = {
			feeds: this.props.feeds,
			lastRefresh: this.props.lastRefresh
		};
		this.getNextFeeds = this.getNextFeeds.bind(this);
		this.onRowClick = this.onRowClick.bind(this);
		this.refresh = this.refresh.bind(this);
	}

	getNextFeeds() {
		return new Promise((resolve, reject) => {
			agent
				.get(this.context.config.rootURL + '/feeds/next')
				.end((err, rsp) => {
					if (err) return reject();
					resolve(rsp.body);
				});
		});
	}

	onRowClick() {
		console.log(arguments);
	}

	refresh() {
		this.setState({
			feeds: [],
			lastRefresh: null
		});
		this.getNextFeeds().then((feeds) => {
			this.setState({
				feeds: feeds,
				lastRefresh: new Date()
			});
		});
	}

	render() {

		const rows = [];
		this.state.feeds.forEach((feed) => {
			rows.push(<FeedItem feed={feed}/>);
		});

		return (
			<div>

				<FeedToolbar lastRefresh={this.state.lastRefresh} onRefreshClick={this.refresh} />

				<Table
					fixedFooter={false}
					fixedHeader={true}
					selectable={true}>
					<TableHeader displaySelectAll={false} enableSelectAll={false} selectAllSelected={false}>
						<FeedItem feed={null} />
					</TableHeader>
					<TableBody
						deselectOnClickAway={false}
						displayRowCheckbox={false}
						onRowSelection={this.onRowClick}
						selectable={true}
						showRowHover={true}
						stripedRows={true}>
						{rows}
					</TableBody>
				</Table>

			</div>
		);

	}

}

export default Feeds;
