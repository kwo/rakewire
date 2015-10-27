import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import test from 'tape';

import Message from '../../app/components/Message.jsx';

test('message rendered with default message', t => {

	const msg = TestUtils.renderIntoDocument(<Message/>);

	const h4Elements = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'h4');
	t.equal(1, h4Elements.length);
	const h4Node = ReactDOM.findDOMNode(h4Elements[0]);
	t.equal('Unconfigured message!', h4Node.innerHTML);

	t.end();

});

test('button should not be rendered', t => {
	const msg = TestUtils.renderIntoDocument(<Message />);
	t.equal(0, TestUtils.scryRenderedDOMComponentsWithTag(msg, 'button').length);
	t.end();
});

test('message rendered with message', t => {

	const text = 'message text';
	const msg = TestUtils.renderIntoDocument(<Message message={text} />);

	const h4 = TestUtils.findRenderedDOMComponentWithTag(msg, 'h4');
	const h4Node = ReactDOM.findDOMNode(h4);
	t.equal(text, h4Node.innerHTML);

	t.end();

});
