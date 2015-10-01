/*

 TODO

 - uglify
 - flatten directory structure

*/

const CleanPlugin = require('clean-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlPlugin = require('html-webpack-plugin');
const path = require('path');
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
			{ test: /\.css$/,  loader: ExtractTextPlugin.extract('style-loader', 'css-loader') },
			{ test: /\.jsx?$/, loader: 'babel?optional[]=runtime&stage=0', exclude: /node_modules/ }
		]
	},
	plugins: [
		new CleanPlugin(['public']),
		new ExtractTextPlugin('styles-[hash].css'),
		new HtmlPlugin({
			template: path.resolve(__dirname, 'app/index.html'),
			inject: 'body'
		}),
		new webpack.optimize.CommonsChunkPlugin('vendor', 'vendor-[hash].js')
	],
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};
