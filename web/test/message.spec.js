import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import should from 'should';
import Message from '../app/components/Message.jsx';

describe('Message', () => {

	it('message rendered with default message', done => {

		const msg = TestUtils.renderIntoDocument(<Message/>);

		const h4 = TestUtils.findRenderedDOMComponentWithTag(msg, 'h4');
		const h4Node = ReactDOM.findDOMNode(h4);
		should.equal('Unconfigured message!', h4Node.innerHTML);

		done();

	});

	it('message rendered with message', done => {

		const text = 'message text';
		const msg = TestUtils.renderIntoDocument(<Message message={text} />);

		const h4 = TestUtils.findRenderedDOMComponentWithTag(msg, 'h4');
		const h4Node = ReactDOM.findDOMNode(h4);
		should.equal(text, h4Node.innerHTML);

		done();

	});

});
