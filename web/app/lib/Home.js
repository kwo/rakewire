import React, { PropTypes } from 'react';

class Home extends React.Component {

	static displayName = 'home';

	// static propTypes = {
	// 	title: PropTypes.string
	// };

	static contextTypes = {
		muiTheme : PropTypes.object.isRequired
	};

	// static defaultProps = {
	// 	title: 'title'
	// }

	constructor(props, context) {
		super(props, context);
		this.state = {};
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
