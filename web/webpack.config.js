/*

 TODO

 - add hash: app-[hash].js
 - modify index.html with asset/hash names
 - split into js and css
 - uglify

*/

const Clean = require('clean-webpack-plugin');

module.exports = {
	entry: './app/lib/index.js',
	output: {
		path: './public',
		filename: 'app.js'
	},
	module: {
		loaders: [
			{ test: /\.css$/,  loader: 'style!css' },
			{ test: /\.jsx?$/, loader: 'babel?optional[]=runtime&stage=0', exclude: /node_modules/ }
		]
	},
	plugins: [
		new Clean(['public'])
	],
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};
