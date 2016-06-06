import React, {PropTypes} from 'react';
import {List, ListItem} from 'material-ui/List';

const style = {

	tree: {
		float: 'left',
		width: '15em',
		marginRight: '1em'
	},

	treelist: {
		paddingTop: 0,
		paddingBottom: 0
	},

	viewer: {
		overflow: 'hidden'
	},

	listitem: {
		lineHeight: '67%',
		fontSize: 'small'
	}

};

class Reader extends React.Component {

	static displayName = 'reader';

	static propTypes = {
		children: PropTypes.node
	};

	static defaultProps = {
	}

	constructor(props, context) {
		super(props, context);
		this.state = {};
	}

	componentDidMount() {
	}

	componentWillUnmount() {
	}

	render() {
		return (
			<div>

				<div id="tree" style={style.tree}>
					<List id="treelist" style={style.treelist}>
						<ListItem primaryText="All" style={style.listitem} />
						<ListItem primaryText="Group1" style={style.listitem} />
						<ListItem primaryText="Group2" style={style.listitem} />
						<ListItem primaryText="Group3" style={style.listitem}  />
						<ListItem primaryText="Group4" style={style.listitem}
							initiallyOpen={false}
							primaryTogglesNestedList={true}
							nestedItems={[
								<ListItem primaryText="Feed1" style={style.listitem} />,
								<ListItem primaryText="Feed2" style={style.listitem} />,
								<ListItem primaryText="Feed3" style={style.listitem} />,
								<ListItem primaryText="Feed4" style={style.listitem} />
							]}/>
					</List>
				</div>

				<div style={style.viewer}>Viewer</div>

			</div>
		);
	}

}

export default Reader;
