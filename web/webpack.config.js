/*

 TODO

 - split into js and css
 - uglify

*/

const Clean = require('clean-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
	entry: './app/lib/index.js',
	output: {
		path: './public',
		filename: 'app-[hash].js'
	},
	module: {
		loaders: [
			{ test: /\.css$/,  loader: 'style!css' },
			{ test: /\.jsx?$/, loader: 'babel?optional[]=runtime&stage=0', exclude: /node_modules/ }
		]
	},
	plugins: [
		new Clean(['public']),
		new HtmlWebpackPlugin({
			template: 'app/index.html',
			inject: 'body'
		})
	],
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};
