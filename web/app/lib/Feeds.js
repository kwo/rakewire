import React, { PropTypes } from 'react';
import {Alert, Button, Table} from 'react-bootstrap';
import FeedEntry from './components/FeedEntry';

class Feeds extends React.Component {

	static displayName = 'feeds';

	static contextTypes = {
		config: PropTypes.object.isRequired
	};

	constructor(props, context) {
		super(props, context);
		this.state = {
			notification: null,
			feeds: null,
			lastRefresh: null
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
		//this.saveState();
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
			notification: {type: 'info', message: 'loading feeds...'},
			feeds: null,
			lastRefresh: null
		});
		this.getNextFeeds()
			.then(feeds => {
				this.setState({
					notification: null,
					feeds: feeds,
					lastRefresh: new Date()
				});
				this.saveState();
			})
			.catch(e => {
				this.setState({
					notification: {type: 'warning', message: `Cannot load feeds: ${e}`},
					feeds: null,
					lastRefresh: null
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

		if (this.state.notification) {
			let refreshButton = '';
			if (this.state.notification.type === 'warning') {
				refreshButton = (
					<p>
						<Button bsStyle="default" onClick={this.refresh}>
							Refresh
						</Button>
					</p>
				);
			}
			return (
				<Alert bsStyle={this.state.notification.type}>
					<h4>{this.state.notification.message}</h4>
					{refreshButton}
				</Alert>
			);
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
