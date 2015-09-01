import 'site.css!'; // #TODO:0 css is missing entries for individual icons

import React from 'react';
import { Styles } from 'material-ui';
import injectTapEventPlugin from 'npm:react-tap-event-plugin@0.1.7/src/injectTapEventPlugin';

const ThemeManager = new Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.LIGHT);
injectTapEventPlugin();

import Router from 'react-router';
import routes from './routes';
Router.run(routes, Router.HistoryLocation, function(Handler) {
	React.render(<Handler />, document.getElementById('app'));
});
