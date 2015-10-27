import React from 'react';
import TestUtils from 'react-addons-test-utils';
import should from 'should';
import Message from '../app/components/Message.jsx';

describe('Message', () => {
	it('message rendered into dom', done => {
		const msg = TestUtils.renderIntoDocument(<Message/>);
		should.exist(msg);
		done();
	});
});
