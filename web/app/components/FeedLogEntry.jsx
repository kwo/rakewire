import React, { PropTypes } from 'react';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';

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

		const style = {
			oh: {
				textAlign: 'center'
			},
			num: {
				textAlign: 'right'
			}
		};

		const logEntry = this.props.logEntry;

		if (!logEntry) {
			return (
				<thead>
					<tr>
						<th colSpan="3" style={style.oh}>General</th>
						<th colSpan="6" style={style.oh}>HTTP</th>
						<th colSpan="6" style={style.oh}>Feed</th>
					</tr>
					<tr>
						<th>Start Time</th>
						<th>Duration</th>
						<th>Result</th>

						<th>Status Code</th>
						<th>Size</th>
						<th>Gzip</th>
						<th>Mime</th>
						<th>ETag</th>
						<th>Last-Modified</th>

						<th>Flavor</th>
						<th>Generator</th>
						<th>Title</th>
						<th>New</th>
						<th>Total</th>
						<th>Updated</th>
					</tr>
				</thead>
			);
		}

		const formatDateTime = function(dt) {
			if (dt && dt.year() <= 1) return '';
			return dt.format('YYYY-MM-DD HH:mm');
		};

		const formatBool = function(value, name) {
			if (value) {
				const elementId = `${logEntry.startTime.format('x')} - ${name}`;
				return (
					<OverlayTrigger overlay={<Tooltip id={elementId}>{String(value)}</Tooltip>} placement="top">
						<i className="material-icons"
							style={{cursor: 'pointer'}}>
							done
						</i>
					</OverlayTrigger>
				);
			}
			return '';
		};

		const formatValue = function(value, message, name) {
			if (message) {
				const elementId = `${logEntry.startTime.format('x')} - ${name}`;
				return (
					<OverlayTrigger overlay={<Tooltip id={elementId}>{String(message)}</Tooltip>} placement="top">
						<span style={{cursor: 'pointer', textDecoration: 'underline'}}>
							{value}
						</span>
					</OverlayTrigger>
				);
			}
			return value;
		};

		const getReadableFileSizeString = function(fileSizeInBytes) {
			if (!fileSizeInBytes) return '';
			let i = -1;
			const byteUnits = ['K', ' MB', ' GB', ' TB', 'PB', 'EB', 'ZB', 'YB'];
			do {
				fileSizeInBytes = fileSizeInBytes / 1024;
				i++;
			} while (fileSizeInBytes > 1024);
			return Math.max(fileSizeInBytes, 0.1).toFixed(0) + byteUnits[i];
		};

		return (
				<tr>
					<td>{formatDateTime(logEntry.startTime)}</td>
					<td style={style.num}>{logEntry.duration / 1000000}</td>
					<td>{formatValue(logEntry.result, logEntry.resultMessage, 'result')}</td>

					<td style={style.num}>{logEntry.statusCode}</td>
					<td style={style.num}>{getReadableFileSizeString(logEntry.contentLength)}</td>
					<td>{formatBool(logEntry.gzip, 'gzip')}</td>
					<td>{logEntry.contentType}</td>
					<td>{formatBool(logEntry.etag, 'etag')}</td>
					<td>{formatBool(logEntry.lastModified, 'lastModified')}</td>

					<td>{logEntry.flavor}</td>
					<td>{logEntry.generator}</td>
					<td>{logEntry.title}</td>
					<td>{logEntry.newEntries}</td>
					<td>{logEntry.entryCount}</td>
					<td>{formatDateTime(logEntry.lastUpdated)}</td>

				</tr>
		);

	}

}

export default FeedLogEntry;
