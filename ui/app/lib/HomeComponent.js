import React, { PropTypes } from 'react';

class Home extends React.Component {

	static displayName = 'home';

	// static propTypes = {
	// 	title: PropTypes.string
	// };

	static contextTypes = {
		router: PropTypes.func,
		muiTheme : PropTypes.object.isRequired
	};

	static childContextTypes = {
		muiTheme : PropTypes.object
	};

	// static defaultProps = {
	// 	title: 'title'
	// }

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	getChildContext () {
		return {
			muiTheme : this.context.muiTheme
		};
	}

	render() {

		return (
			<div>
				<p>Welcome to Rakewire.</p>
			</div>
		);

	}

}

export default Home;
