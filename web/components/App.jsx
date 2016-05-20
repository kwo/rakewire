import React, {PropTypes} from 'react';
import TopBar from './TopBar';

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		children: PropTypes.node
	};

	constructor(props) {
		super(props);
	}

	render() {
		return (
			<div>
				<TopBar/>
				{this.props.children}
			</div>
		);
	}

}

export default App;
