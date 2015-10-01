import React, {PropTypes}  from 'react';
import {Link} from 'react-router';

class NavLink extends React.Component {

	static displayName = 'navlink';

	static propTypes = {
		children: PropTypes.node,
		to: PropTypes.string.isRequired
	};

	static defaultProps = {
		to: 'home'
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {
		return (
			<li role="presentation"><Link to={this.props.to}>{this.props.children}</Link></li>
		);
	} // render

}

export default NavLink;
