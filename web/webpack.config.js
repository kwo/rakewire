const CleanPlugin = require('clean-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlPlugin = require('html-webpack-plugin');
const path = require('path');
const webpack = require('webpack');

const config = {
	entry: {
		app: path.resolve(__dirname, 'app/index.js'),
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
		new webpack.optimize.CommonsChunkPlugin('vendor', 'vendor-[hash].js'),
		new webpack.optimize.UglifyJsPlugin({
			compress: {
				warnings: false
			},
			mangle: {
				except: ['$super', '$', 'exports', 'require']
			},
			sourceMap: false
		})
	],
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};

const debugMode = function() {
	/* eslint no-var: 0 */
	for (var i = 0; i < process.argv.length; i++) {
		if (process.argv[i] == '--debug') return true;
	}
	return false;
}();

if (debugMode) {
	config.plugins.pop(); // remove uglify to speed up process while developing
}

module.exports = config;
