require('babel-core/register')({
	optional: [ 'es7.classProperties' ]
});

const jsdom = require('jsdom');
const doc = jsdom.jsdom('<!doctype html><html><body></body></html>')
const win = doc.defaultView

// set globals for mocha that make access to document and window feel natural in the test environment
global.document = doc
global.window = win
global.navigator = win.navigator;
