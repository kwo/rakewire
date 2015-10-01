/*

 TODO
 - add hash: app-[hash].js
 - modify index.html with asset/hash names
 - split into js and css
 - uglify
 - clean old assets

*/

module.exports = {
	entry: './app/lib/main.jsx',
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
	resolve: {
		extensions: ['', '.js', '.json', '.jsx']
	}
};
