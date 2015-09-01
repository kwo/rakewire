import 'site.css!';
import React from 'react';
import Router from 'react-router';
import { Styles } from 'material-ui';
import injectTapEventPlugin from 'npm:react-tap-event-plugin@0.1.7/src/injectTapEventPlugin';
import routes from './routes';

const ThemeManager = new Styles.ThemeManager();
ThemeManager.setTheme(ThemeManager.types.LIGHT);
injectTapEventPlugin();

Router.run(routes, Router.HistoryLocation, function(Handler) {
	React.render(<Handler />, document.getElementById('app'));
});
