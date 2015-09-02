/* eslint no-unused-vars: 0 */
import System from 'systemjs';
import '../config.js';
import test from 'tape';

import React from 'react/addons';
const TestUtils = React.addons.TestUtils;


test('Environment', (t) => {
	t.ok(React);
	t.ok(TestUtils);
	t.end();
});

test('Title', (t) => {
	const TitleElement = require('../lib/components/Title');
	t.ok(TitleElement);
	t.end();
});
