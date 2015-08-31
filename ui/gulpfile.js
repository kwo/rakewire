(function() {

	'use strict';

	const gulp = require('gulp');
	const gutil = require('gulp/node_modules/gulp-util');
	const paths = makePaths();

	gulp.task('default', ['build']);
	gulp.task('lint', lint);
	gulp.task('clean', clean);
	gulp.task('resources', ['clean'], resources);
	gulp.task('build', ['lint', 'resources'], build);
	gulp.task('webapp', webapp);

	function makePaths() {

		let src = 'app';
		let dst = 'build';

		let p = {
			src: {
				base: src,
				all:  src + '/**/*',
				js:   src + '/lib/**/*.js'
			},
			dst: {
				base: dst,
				all:  dst + '/**'
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
		log('cleaning...');
		const del = require('del');
		return del(paths.dst.all, {dot: true});
	}

	function resources() {
		log('copying resources...');
		const path = require('path');
		const htmlreplace = require('gulp-html-replace');
		let promises = [];

		promises.push(new Promise(function(resolve, reject) {
			gulp.src(path.join(paths.src.base, 'index.html'))
				.pipe(htmlreplace({
					'js':  'app.js',
					'css': 'app.css'
				}))
				.pipe(gulp.dest(paths.dst.base))
				.on('end', resolve)
				.on('error', reject);
		}));

		promises.push(new Promise(function(resolve, reject) {
			gulp.src(path.join(paths.src.base, '*.txt'))
				.pipe(gulp.dest(paths.dst.base))
				.on('end', resolve)
				.on('error', reject);
		}));

		return Promise.all(promises);

	}

	function build() {
		log('building...');
		const jspm = require('jspm');
		const path = require('path');
		return jspm.bundleSFX('lib/main', path.join(paths.dst.base, 'app.js'),
			{ mangle: true, minify: true, lowResSourceMaps: false, sourceMaps: false }
		);
	}

	function webapp() {

		let src = 'build';
		let mode = process.argv.slice(3);

		if (mode == '-app') {
			src = 'app';
		} else if (mode == '-build') {
			src = 'build';
		} else {
			return log('unknown mode: ' + mode);
		}

		const symlink = require('gulp-sym');

		gulp
			.src(src)
			.pipe(symlink(function(/*source*/) {
				return 'webapp';
			}, { force: true, relative: true }));

	}

})();
