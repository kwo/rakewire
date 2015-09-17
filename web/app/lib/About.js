import React, { PropTypes } from 'react';

class About extends React.Component {

	static displayName = 'about';

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
				<p>About Rakewire</p>
			</div>
		);

	}

}

export default About;
