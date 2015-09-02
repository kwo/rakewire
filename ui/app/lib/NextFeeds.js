import React, { PropTypes } from 'react';
import { FlatButton, FontIcon } from 'material-ui';
import agent from 'superagent';
import FeedRow from './components/FeedRow';

// #DOING:10 lose state after route change - save to app-wide repository - localstate
// #DOING:10 save last refresh time, auto-reload if state too old

class NextFeeds extends React.Component {

	static displayName = 'nextfeeds';

	static propTypes = {
		feeds: PropTypes.array
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
			feeds: this.props.feeds
		};
		this.getNextFeeds = this.getNextFeeds.bind(this);
		this.onRefreshClick = this.onRefreshClick.bind(this);
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

	onRefreshClick() {
		this.getNextFeeds().then((feeds) => {
			this.setState({
				feeds: feeds
			});
			console.log(feeds);
		});
	}

	render() {

		const rows = [];
		this.state.feeds.forEach((feed) => {
			rows.push(<FeedRow feed={feed}/>);
		});

		return (
			<div>

				<FlatButton
					label="Refresh"
					onTouchTap={this.onRefreshClick}
					secondary={true}>
					<FontIcon className="material-icons">refresh</FontIcon>
				</FlatButton>

				<hr/>

				<table>

					<thead>
						<tr>
							<th>Next</th>
							<th>Last</th>
							<th>Status</th>
							<th>Code</th>
							<th>Updated</th>
							<th>Check</th>
							<th>Feed</th>
						</tr>
					</thead>

					<tbody>
						{rows}
					</tbody>

				</table>

			</div>
		);

	}

}

export default NextFeeds;
