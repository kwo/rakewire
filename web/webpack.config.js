const webpack = require('webpack');
const path = require('path');
const CleanPlugin = require('clean-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlPlugin = require('html-webpack-plugin');

const config = {
	entry: {
		app: path.resolve(__dirname, 'app.js'),
		vendor: [
			'moment',
			'react',
			'react-dom',
			'react-router',
			'react-tap-event-plugin',
			'whatwg-fetch'
		]
	},
	output: {
		path: path.resolve(__dirname, 'public'),
		filename: 'app.js'
	},
	module: {
		loaders: [
			{ test: /\.(css|scss)$/,  loader: ExtractTextPlugin.extract('style', 'css!sass') },
			{ test: /\.jsx?$/,
				loader: 'babel',
				exclude: /node_modules/,
				query: {
					cacheDirectory: true,
					presets: ['es2015', 'react', 'stage-0']
				}
			},
			{ test: /\.jsx?$/, loader: 'eslint-loader', exclude: /node_modules/ },
			{ test: /\.(eot|svg|ttf|woff|woff2)$/, loader: 'file?name=fonts/[name].[ext]' }
		]
	},
	plugins: [
		new CleanPlugin(['public']),
		new ExtractTextPlugin('app.css'),
		new HtmlPlugin({
			template: path.resolve(__dirname, 'index.html'),
			inject: 'body'
		}),
		new webpack.optimize.CommonsChunkPlugin('vendor', 'vendor.js'),
		new webpack.optimize.UglifyJsPlugin({
			compress: {
				warnings: false
			},
			mangle: {
				except: ['$super', '$', 'exports', 'require']
			},
			sourceMap: false
		}),
		new webpack.DefinePlugin({
			'process.env': {
				'NODE_ENV': '"production"'
			}
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
	config.plugins.pop(); // remove NODE_ENV production
	config.plugins.pop(); // remove uglify to speed up process while developing
}

module.exports = config;
