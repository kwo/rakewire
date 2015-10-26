import React, {PropTypes} from 'react';
import AppHeader from '../components/AppHeader';

const Config = {
	rootURL: '/api'
};

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		children: PropTypes.node,
		title: PropTypes.string
	};

	static childContextTypes = {
		config: PropTypes.object
	};

	static defaultProps = {
		title: 'Rakewire'
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	getChildContext () {
		return {
			config: Config
		};
	}

	render() {

		return (
			<div>
				<AppHeader />
				<div className="container-fluid">
					{this.props.children}
				</div>
			</div>
		);

	}

}

export default App;
