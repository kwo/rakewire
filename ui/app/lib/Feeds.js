import React, { PropTypes } from 'react';
import { Table, TableBody, TableHeader } from 'material-ui';
import FeedItem from './components/FeedItem';
import FeedToolbar from './components/FeedToolbar';

// XXX lose state after route change - save to app-wide repository - localstate
// XXX auto-reload if state too old

class Feeds extends React.Component {

	static displayName = 'feeds';

	static contextTypes = {
		config: PropTypes.object.isRequired,
		muiTheme: PropTypes.object.isRequired
	};

	constructor(props, context) {
		super(props, context);
		this.state = {
			feeds: [],
			lastRefresh: null
		};
		this.getNextFeeds = this.getNextFeeds.bind(this);
		this.refresh = this.refresh.bind(this);
	}

	componentDidMount() {
		this.refresh(); // XXX: load from local state somehow
		// XXX: live refresh via server-sent events (needs client-side update)
	}

	getNextFeeds() {
		return new Promise((resolve, reject) => {
			fetch(this.context.config.rootURL + '/feeds/next')
				.then(rsp => rsp.json())
				.then(feeds => resolve(feeds))
				.catch(e => reject(e));
		});
	}

	refresh() {
		this.setState({
			feeds: [],
			lastRefresh: null
		});
		this.getNextFeeds()
			.then(feeds => {
				this.setState({
					feeds: feeds,
					lastRefresh: new Date()
				});
			})
			.catch(e => console.log('Cannot load feeds:', e)); // XXX: display on view
	}

	render() {

		const rows = [];
		this.state.feeds.forEach((feed) => {
			rows.push(<FeedItem feed={feed} key={feed.id} />);
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
