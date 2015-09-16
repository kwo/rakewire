import React, { PropTypes } from 'react';
import FeedInfo from './components/FeedInfo';

class Feed extends React.Component {

	static displayName = 'feed';

	static propTypes = {
		params: PropTypes.object.isRequired
	};

	static contextTypes = {
		config: PropTypes.object.isRequired,
		muiTheme: PropTypes.object
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
				const feed = results[0].json();
				feed.log = results[1].json();
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

		return (

			<table style={style.table}>

				<FeedInfo feed={feed}/>

			</table>

		);

	}

}

export default Feed;
