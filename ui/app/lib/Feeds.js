import React, { PropTypes } from 'react';
import agent from 'superagent';
import FeedItem from './components/FeedItem';
import FeedToolbar from './components/FeedToolbar';

// #DOING:50 lose state after route change - save to app-wide repository - localstate
// #DOING:60 auto-reload if state too old

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
		feeds: []
	}

	constructor(props, context) {
		super(props, context);
		this.state = {
			feeds: this.props.feeds,
			lastRefresh: this.props.lastRefresh
		};
		this.getNextFeeds = this.getNextFeeds.bind(this);
		this.refresh = this.refresh.bind(this);
		if (!this.state.lastRefresh) this.refresh();
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

				<table>
					<thead>
						<FeedItem feed={null} />
					</thead>
					<tbody>
						{rows}
					</tbody>
				</table>

			</div>
		);

	}

}

export default Feeds;
