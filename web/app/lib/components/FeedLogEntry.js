import React, { PropTypes } from 'react';
import moment from 'moment';

class FeedLogEntry extends React.Component {

	static displayName = 'feedlogentry';

	static propTypes = {
		logEntry: PropTypes.object
	};

	static defaultProps = {
		logEntry: null
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		const logEntry = this.props.logEntry;

		if (!logEntry) {
			return (
				<tr>
					<td>Start Time</td>
				</tr>
			);
		}

		const formatDateTime = function(dt) {
			return moment(dt).format('YYYY-MM-DD HH:mm');
		};

		return (
			<tr>
				<td>{formatDateTime(logEntry.startTime)}</td>
			</tr>
		);

	}

}

export default FeedLogEntry;
