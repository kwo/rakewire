const uuid = require('node-uuid');
const buildversion = uuid.v1().substr(0, 8);

module.exports = {
	entry: './app/lib/main.jsx',
	output: {
		path: './build',
		filename: `app-${buildversion}.js`
	},
	module: {
		loaders: [
			{ test: /\.css$/,  loader: 'style-loader!css-loader' },
			{ test: /\.jsx?$/, loader: 'babel?optional[]=runtime&stage=0', exclude: /node_modules/ }
		]
	},
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};
