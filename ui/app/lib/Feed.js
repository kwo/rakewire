import React, { PropTypes } from 'react';
//import moment from 'moment';

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
			},
			td: {
				height: 24,
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			},
			th: {
				paddingLeft: 12,
				paddingRight: 12,
				textAlign: 'left',
				width: 40
			}
		};

		const feed = this.state.feed;
		const title = feed.title || feed.last200.feed.title || feed.url;

		return (

			<table style={style.table}>

				<thead>

					<tr key={0}>
						<th style={style.th}>Title</th>
						<td style={style.td}>{title}</td>
					</tr>

					<tr key={1}>
						<th style={style.th}>URL</th>
						<td style={style.td}>{feed.url}</td>
					</tr>

				</thead>

			</table>

		);

	}

}

export default Feed;
