import test from 'tape';

test('Test 1', (assert) => {
	assert.pass('This test will pass.');
	assert.end();
});

test('Test 2', (assert) => {
	const expected = 'something to test';
	const actual = 'sonething to test';
	assert.equal(actual, expected);
	assert.end();
});
