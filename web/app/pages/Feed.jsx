import React, { PropTypes } from 'react';
import {Table} from 'react-bootstrap';
import moment from 'moment';
import FeedInfo from '../components/FeedInfo';
import FeedLogEntry from '../components/FeedLogEntry';
import Message from '../components/Message';

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
			message: null,
			feed: null
		};
		this.refresh = this.refresh.bind(this);
	}

	componentDidMount() {
		this.refresh();
	}

	refresh() {
		// no message at start -- too much flicker
		if (this.props.params.id) {
			const feedURL = `${this.context.config.rootURL}/feeds/${this.props.params.id}`;
			const feedLogURL = `${feedURL}/log`;
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
					feed: feed,
					message: null
				});
			})
			.catch(e => {
				this.setState({
					feed: null,
					message: {type: 'warning', text: `Cannot load feed: ${e}`}
				});
			});
		} // id
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
