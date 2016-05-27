import React, {PropTypes} from 'react';
import TopBar from './TopBar';

const style = {
	root: {
		fontFamily: 'Roboto, sans-serif'
	},
	content: {
		margin: '1em'
	}
};

class App extends React.Component {

	static displayName = 'app';

	static propTypes = {
		children: PropTypes.node,
		sidebar: PropTypes.bool,
		title: PropTypes.string
	};

	static defaultProps = {
		sidebar: false,
		title: 'Rakewire'
	}

	// static childContextTypes = {
	// };

	constructor(props, ctx) {
		super(props, ctx);
		this.state = {};
	}

	render() {
		return (
			<div style={style.root}>
				<TopBar sidebar={this.props.sidebar} title={this.props.title} />
				<div style={style.content}>
					{this.props.children}
				</div>
			</div>
		);
	}

}

export default App;
