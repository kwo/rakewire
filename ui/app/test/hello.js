import test from 'tape';
import React from 'react/addons';
const TestUtils = React.addons.TestUtils;


test('Environment', (t) => {
	t.ok(React);
	t.ok(TestUtils);
	t.end();
});

test('Title', (t) => {
	const TitleElement = require('./components/Title');
	t.ok(TitleElement);
	t.end();
});
