import React, {PropTypes} from 'react';

class Icon extends React.Component {

	static displayName = 'icon';

	static propTypes = {
		name: PropTypes.string.isRequired,
		style: PropTypes.object
	};

	constructor(props) {
		super(props);
	}

	render() {
		return (
			<i	className="material-icons" style={Object.assign({ lineHeight: 'inherit' }, this.props.style)}>{this.props.name}</i>
		);
	}

}

export default Icon;
