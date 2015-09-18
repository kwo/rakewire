import React, { PropTypes } from 'react';
import {Table} from 'react-bootstrap';
import FeedInfo from './components/FeedInfo';
import FeedLogEntry from './components/FeedLogEntry';
import moment from 'moment';

class Feed extends React.Component {

	static displayName = 'feed';

	static propTypes = {
		params: PropTypes.object.isRequired
	};

	static contextTypes = {
		config: PropTypes.object.isRequired
	};

	static defaultProps = {
		params: {}
	}

	constructor(props, context) {
		super(props, context);
		this.state = {
			feed: null
		};
	}

	componentDidMount() {
		if (this.props.params.id) {
			const feedURL = this.context.config.rootURL + '/feeds/' + this.props.params.id;
			const feedLogURL = feedURL + '/log';
			Promise.all([
				fetch(feedURL),
				fetch(feedLogURL)
			])
			.then(results => {
				const p = [];
				results.forEach(r => {
					p.push(r.json());
				});
				return Promise.all(p);
			})
			.then(results => {
				const feed = results[0];
				feed.log = results[1];
				feed.log.forEach(log => {
					log.startTime = moment(log.startTime);
					log.feed.updated = moment(log.feed.updated);
				});
				return feed;
			})
			.then(feed => {
				this.setState({
					feed: feed
				});
			})
			.catch(e => console.error(e)); // XXX: display error in UI
		} // id
	}

	render() {

		if (!this.state.feed) {
			return (
				<div>loading...</div>
			);
		}

		const feed = this.state.feed;

		const logEntries = [];
		feed.log.forEach(logEntry => {
			logEntries.push(<FeedLogEntry key={logEntry.startTime.format('x')} logEntry={logEntry} />);
		});

		return (

			<div>

				<Table condensed={true} hover={true} responsive={true}>
					<FeedInfo feed={feed}/>
				</Table>

				<Table>
					(<FeedLogEntry logEntry={null} />
					<tbody>
						{logEntries}
					</tbody>
				</Table>

			</div>

		);

	}

}

export default Feed;
