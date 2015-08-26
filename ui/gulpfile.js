(function() {

	'use strict';

	const gulp = require('gulp');
	const gutil = require('gulp/node_modules/gulp-util');
	const paths = makePaths();

	gulp.task('default', ['build']);
	gulp.task('lint', lint);
	gulp.task('clean', clean);
	gulp.task('build', ['clean', 'lint'], build);
	gulp.task('deploy', ['build'], deploy);

	function makePaths() {

		let base = 'app';
		let src = base + '/lib';
		let dst = base;

		let p = {
			dst: {
				base: dst,
				all:  dst + '/**/*'
			},
			src: {
				base: src,
				all:  src + '/**/*',
				js:   src + '/**/*.js'
			}
		};

		return p;

	}

	function log() {
		gutil.log.apply(null, arguments);
	}

	function lint() {
		log('linting...');
		const eslint = require('gulp-eslint');
		return gulp.src(paths.src.js)
			.pipe(eslint())
			.pipe(eslint.format());
	}

	function clean() {
		log('cleaning... temporairly disabled');
		// const del = require('del');
		// return new Promise(function(resolve, reject) {
		// 	del([paths.dst.all, paths.dst.base], {dot: true}, function(err) {
		// 		if (err) return reject(err);
		// 		resolve();
		// 	});
		// });
	}

	function build() {
		log('building...');
		const jspm = require('jspm');
		const path = require('path');
		return jspm.bundle(
			'lib/main',
			path.join(paths.dst.base, 'bundle.js'),
			{ inject: true, mangle: true, minify: true, lowResSourceMaps: false, sourceMaps: false }
		);
	}

	function deploy() {
		log('not implemented');
	}

})();
