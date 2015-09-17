import React, { PropTypes } from 'react';
import moment from 'moment';

class FeedLogEntry extends React.Component {

	static displayName = 'feedlogentry';

	static propTypes = {
		logEntry: PropTypes.object
	};

	static contextTypes = {
		muiTheme : PropTypes.object
	};

	static defaultProps = {
		logEntry: null
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	render() {

		const style = {
			th: {
				paddingLeft: '6px',
				paddingRight: '6px',
				textAlign: 'left',
				width: '40px'
			},
			td: {
				height: '24px',
				paddingLeft: '6px',
				paddingRight: '6px',
				textAlign: 'left',
				width: '40px'
			},
		};

		const logEntry = this.props.logEntry;

		if (!logEntry) {
			return (
				<tr>
					<th style={style.th}>Start Time</th>
				</tr>
			);
		}

		const formatDateTime = function(dt) {
			return moment(dt).format('YYYY-MM-DD HH:mm');
		};

		return (
			<tr>
				<td style={style.td}>{formatDateTime(logEntry.startTime)}</td>
			</tr>
		);

	}

}

export default FeedLogEntry;
