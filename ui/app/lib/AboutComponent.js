import React, { PropTypes } from 'react';

class About extends React.Component {

	static displayName = 'about';

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
				<p>About Rakewire</p>
			</div>
		);

	}

}

export default About;
