import React, {PropTypes} from 'react';

class MyComponent extends React.Component {

	static displayName = 'mycomponent';

	static propTypes = {
		children: PropTypes.node
	};

	static defaultProps = {
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	componentDidMount() {
	}

	componentWillUnmount() {
	}

	render() {
		return (
			<div>
				<p>About Rakewire</p>
			</div>
		);
	}

}

export default MyComponent;
