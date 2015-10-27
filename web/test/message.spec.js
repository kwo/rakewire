import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import {assert} from 'chai';
import Message from '../app/components/Message.jsx';

describe('Message', () => {

	it('message rendered with default message', done => {

		const msg = TestUtils.renderIntoDocument(<Message/>);

		const h4Elements = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'h4');
		assert.strictEqual(1, h4Elements.length);
		const h4Node = ReactDOM.findDOMNode(h4Elements[0]);
		assert.strictEqual('Unconfigured message!', h4Node.innerHTML);

		done();

	});

	it('button should not be rendered', done => {
		const msg = TestUtils.renderIntoDocument(<Message />);
		assert.strictEqual(0, TestUtils.scryRenderedDOMComponentsWithTag(msg, 'button').length);
		done();
	});

	it('message rendered with message', done => {

		const text = 'message text';
		const msg = TestUtils.renderIntoDocument(<Message message={text} />);

		const h4 = TestUtils.findRenderedDOMComponentWithTag(msg, 'h4');
		const h4Node = ReactDOM.findDOMNode(h4);
		assert.strictEqual(text, h4Node.innerHTML);

		done();

	});

});
