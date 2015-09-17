import React, { PropTypes } from 'react';
import FeedEntry from './components/FeedEntry';
import FeedToolbar from './components/FeedToolbar';

class Feeds extends React.Component {

	static displayName = 'feeds';

	static contextTypes = {
		config: PropTypes.object.isRequired
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
		// TODO: slim size of data, perhaps change api for slimmer payload or use paging

		const asyncParse = function(json) {
			if (!json) return Promise.reject('Skipping, no json provided.');
			return (new Response(json)).json();
		};

		asyncParse(sessionStorage.getItem('feeds.state'))
			.then(state => this.setState(state))
			.catch(e => console.log('Cannot load state', e)); // TODO: handle errors in UI

	}

	componentWillUnmount() {
		sessionStorage.setItem('feeds.state', JSON.stringify(this.state));
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
				sessionStorage.setItem('feeds.state', JSON.stringify(this.state));
			})
			.catch(e => console.log('Cannot load feeds:', e)); // XXX: display error in UI
	}

	render() {

		const rows = [];
		this.state.feeds.forEach((feed) => {
			rows.push(<FeedEntry feed={feed} key={feed.id} />);
		});

		return (
			<div>

				<FeedToolbar lastRefresh={this.state.lastRefresh} onRefreshClick={this.refresh} />

				<table>
					<thead>
						<FeedEntry feed={null} />
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
