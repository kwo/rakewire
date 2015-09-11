import React, { PropTypes } from 'react';

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
			fetch(this.context.config.rootURL + '/feeds/' + this.props.params.id)
				.then(rsp => rsp.json())
				.then(feed => {
					this.setState({
						feed: feed
					});
				})
				.catch(e => {
					console.error(e); // XXX: display error in UI
				});
		} // id
	}

	render() {

		if (!this.state.feed) {
			return (
				<div>loading...</div>
			);
		}

		return (
			<div>
				{this.state.feed.id}: {this.state.feed.title}
			</div>
		);

	}

}

export default Feed;
