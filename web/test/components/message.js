import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';
import test from 'tape';

import Message from '../../app/components/Message.jsx';

test('message rendered with default message', t => {

	const msg = TestUtils.renderIntoDocument(<Message/>);

	const components = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'h4');
	t.equal(components.length, 1);
	const element = ReactDOM.findDOMNode(components[0]);
	t.equal(element.innerHTML, 'Unconfigured message!');

	t.end();

});

test('message rendered with message', t => {

	const text = 'message text';
	const msg = TestUtils.renderIntoDocument(<Message message={text} />);

	const component = TestUtils.findRenderedDOMComponentWithTag(msg, 'h4');
	const element = ReactDOM.findDOMNode(component);
	t.equal(element.innerHTML, text);

	t.end();

});

test('default message type is info', t => {

	const msg = TestUtils.renderIntoDocument(<Message/>);

	const components = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'div');
	t.equal(components.length, 1);
	const element = ReactDOM.findDOMNode(components[0]);
	t.equal(element.getAttribute('class'), 'alert alert-info');

	t.end();

});

test('button should not be rendered when no props', t => {
	const msg = TestUtils.renderIntoDocument(<Message/>);
	t.equal(TestUtils.scryRenderedDOMComponentsWithTag(msg, 'button').length, 0);
	t.end();
});

test('button should not be rendered when btnLabel but no btnClick', t => {
	const msg = TestUtils.renderIntoDocument(<Message btnLabel={"click me"}/>);
	t.equal(TestUtils.scryRenderedDOMComponentsWithTag(msg, 'button').length, 0);
	t.end();
});

test('button click', t => {

	t.plan(7);

	const buttonClick = function(e) {
		t.ok(e);
	};
	const buttonText = 'click me';
	const messageText = 'hello world';
	const messageType = 'warning';

	const msg = TestUtils.renderIntoDocument(<Message btnClick={buttonClick} btnLabel={buttonText} message={messageText} type={messageType} />);

	// test message text
	const h4Components = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'h4');
	t.equal(h4Components.length, 1);
	const h4Element = ReactDOM.findDOMNode(h4Components[0]);
	t.equal(h4Element.innerHTML, messageText);

	// test message type
	const divComponents = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'div');
	t.equal(divComponents.length, 1);
	const divElement = ReactDOM.findDOMNode(divComponents[0]);
	t.equal(divElement.getAttribute('class'), 'alert alert-warning');

	// test button text & click
	const buttonComponents = TestUtils.scryRenderedDOMComponentsWithTag(msg, 'button');
	t.equal(buttonComponents.length, 1);
	const buttonElement = ReactDOM.findDOMNode(buttonComponents[0]);
	t.equal(buttonElement.innerHTML, buttonText);

	TestUtils.Simulate.click(buttonElement);

});
