import React, { PropTypes } from 'react';

import mui from 'material-ui';
const AutoPrefix = mui.Styles.AutoPrefix;
const ImmutabilityHelper = mui.Utils.ImmutabilityHelper;
const Typography = mui.Styles.Typography;

class Title extends React.Component {

	static displayName = 'title';

	static propTypes = {
		onClick: PropTypes.func.isRequired,
		title: PropTypes.string.isRequired
	};

	static contextTypes = {
		muiTheme : PropTypes.object.isRequired
	};

	static childContextTypes = {
		muiTheme : PropTypes.object
	};

	static defaultProps = {
		title: 'title'
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	getChildContext () {
		return {
			muiTheme : this.context.muiTheme
		};
	}

	// from AppBar
	getStyles() {
		let themeVariables = this.context.muiTheme.component.appBar;
		let styles = {
			title: {
				whiteSpace: 'nowrap',
				overflow: 'hidden',
				textOverflow: 'ellipsis',
				margin: 0,
				paddingTop: 0,
				letterSpacing: 0,
				fontSize: 24,
				fontWeight: Typography.fontWeightNormal,
				color: themeVariables.textColor,
				lineHeight: themeVariables.height + 'px',
			},
			mainElement: {
				boxFlex: 1,
				flex: '1',
			},
			link: {
				cusror: 'pointer'
			}
		};

		return styles;
	}

	mergeAndPrefix() {
		let mergedStyles = ImmutabilityHelper.merge.apply(this, arguments);
		return AutoPrefix.all(mergedStyles);
	}

	render() {

		const appBoxStyles = this.getStyles();
		const titleStyle = this.mergeAndPrefix(appBoxStyles.title, appBoxStyles.mainElement);
		titleStyle.cursor = 'pointer';

		return (
			<h1
				onTouchTap={this.props.onClick}
				style={titleStyle}>
				{this.props.title}
			</h1>
		);
	}

}

export default Title;
