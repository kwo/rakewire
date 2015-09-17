import React, { PropTypes } from 'react';
import moment from 'moment';
import FeedInfo from './components/FeedInfo';
import FeedLogEntry from './components/FeedLogEntry';

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

		const style = {
			table: {
				width: '100%'
			}
		};

		const feed = this.state.feed;

		const logEntries = [];
		feed.log.forEach(logEntry => {
			logEntries.push(<FeedLogEntry key={logEntry.startTime.format('x')} logEntry={logEntry} />);
		});

		return (

			<table style={style.table}>

				<FeedInfo feed={feed}/>

				<tbody>
					{logEntries}
				</tbody>

			</table>

		);

	}

}

export default Feed;
