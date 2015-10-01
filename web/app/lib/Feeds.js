import React, { PropTypes } from 'react';
import {Table} from 'react-bootstrap';
import FeedEntry from './components/FeedEntry';
import Message from './components/Message';

class Feeds extends React.Component {

	static displayName = 'feeds';

	static contextTypes = {
		config: PropTypes.object.isRequired
	};

	constructor(props, context) {
		super(props, context);
		this.state = {
			feeds: null,
			lastRefresh: null,
			message: null
		};
		this.getNextFeeds = this.getNextFeeds.bind(this);
		this.refresh = this.refresh.bind(this);
		this.loadState = this.loadState.bind(this);
		this.saveState = this.saveState.bind(this);
	}

	componentDidMount() {

		const asyncParse = function(json) {
			if (!json) return Promise.reject('Skipping, no json provided.');
			return (new Response(json)).json();
		};

		asyncParse(this.loadState())
			.then(state => this.setState(state))
			.catch(e => {
				console.log('Cannot load feeds.state from localstorage, refreshing.', e);
				this.refresh();
			});

	}

	componentWillUnmount() {
		this.saveState();
	}

	getNextFeeds() {
		return new Promise((resolve, reject) => {
			fetch(`${this.context.config.rootURL}/feeds/next`)
				.then(rsp => rsp.json())
				.then(feeds => resolve(feeds))
				.catch(e => reject(e));
		});
	}

	refresh() {
		this.setState({
			feeds: null,
			lastRefresh: null,
			message: {type: 'info', text: 'loading feeds...'}
		});
		this.getNextFeeds()
			.then(feeds => {
				this.setState({
					feeds: feeds,
					lastRefresh: new Date(),
					message: null
				});
				this.saveState();
			})
			.catch(e => {
				this.setState({
					feeds: null,
					lastRefresh: null,
					message: {type: 'warning', text: `Cannot load feeds: ${e}`}
				});
			});
	}

	loadState() {
		const state = sessionStorage.getItem('feeds.state');
		if (state.lastRefresh) state.lastRefresh = new Date(state.lastRefresh);
		return state;
	}

	saveState() {
		sessionStorage.setItem('feeds.state', JSON.stringify({
			feeds: this.state.feeds,
			lastRefresh: (this.state.lastRefresh) ? this.state.lastRefresh.getTime() : null
		}));
	}

	render() {

		if (this.state.message) {
			const n = this.state.message;
			if (n.type === 'warning') {
				return ( <Message btnClick={this.refresh} btnLabel={"Refresh"} message={n.text} type={n.type} /> );
			} else {
				return ( <Message message={n.text} type={n.type} /> );
			}
		}

		const rows = [];
		if (this.state.feeds) {
			this.state.feeds.forEach((feed) => {
				rows.push(<FeedEntry feed={feed} key={feed.id} />);
			});
		}

		return (
			<div>

				<Table condensed={true} hover={true} responsive={true} striped={true} >

					<thead>
						<FeedEntry feed={null} lastRefresh={this.state.lastRefresh} onRefreshClick={this.refresh} />
					</thead>

					<tbody>
						{rows}
					</tbody>

				</Table>

			</div>
		);

	}

}

export default Feeds;
