import React, {PropTypes}  from 'react';
import {Alert, Button} from 'react-bootstrap';

class Message extends React.Component {

	static displayName = 'message';

	static propTypes = {
		btnClick: PropTypes.func,
		btnLabel: PropTypes.string,
		message: PropTypes.string.isRequired,
		type: PropTypes.string,
	};

	static defaultProps = {
		btnLabel: 'OK',
		message: 'Unconfigured message!',
		type: 'info'
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		let refreshButton = '';
		if (this.props.btnClick) {
			refreshButton = (
				<p>
					<Button bsStyle="default" onClick={this.props.btnClick}>
						{this.props.btnLabel}
					</Button>
				</p>
			);
		}

		return (
			<Alert bsStyle={this.props.type}>
				<h4>{this.props.message}</h4>
				{refreshButton}
			</Alert>
		);
	} // render

}

export default Message;
