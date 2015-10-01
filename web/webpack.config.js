/*

 TODO

 - split into js and css
 - uglify
 - flatten directory structure

*/

const path = require('path');
const Clean = require('clean-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

module.exports = {
	entry: {
		app: path.resolve(__dirname, 'app/lib/index.js'),
		vendor: ['moment', 'react', 'react-bootstrap', 'react-router', 'whatwg-fetch']
	},
	output: {
		path: path.resolve(__dirname, 'public'),
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
			template: path.resolve(__dirname, 'app/index.html'),
			inject: 'body'
		}),
		new webpack.optimize.CommonsChunkPlugin('vendor', 'vendor-[hash].js')
	],
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};
