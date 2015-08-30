(function() {

	'use strict';

	const gulp = require('gulp');
	const gutil = require('gulp/node_modules/gulp-util');
	const paths = makePaths();

	gulp.task('default', ['build']);
	gulp.task('lint', lint);
	gulp.task('clean', clean);
	gulp.task('resources', ['clean'], resources);
	gulp.task('build', ['resources'], build);
	gulp.task('deploy', ['build'], deploy);

	function makePaths() {

		let src = 'development';
		let dst = 'production';

		let p = {
			src: {
				base: src,
				all:  src + '/**/*',
				js:   src + '/**/*.js'
			},
			dst: {
				base: dst,
				all:  dst + '/**/*'
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
		return new Promise(function(resolve, reject) {
			del([paths.dst.all, paths.dst.base], {dot: true}, function(err) {
				if (err) return reject(err);
				resolve();
			});
		});
	}

	function resources() {
		log('copying resources...');
		const htmlreplace = require('gulp-html-replace');
		let promises = [];

		promises.push(new Promise(function(resolve, reject) {
			gulp.src(paths.src.base + '/index.html')
				.pipe(htmlreplace({
					'js': 'app.js'
				}))
				.pipe(gulp.dest(paths.dst.base))
				.on('end', resolve)
				.on('error', reject);
		}));

		promises.push(new Promise(function(resolve, reject) {
			gulp.src(paths.src.base + '/site.css')
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

	function deploy() {
		log('not implemented');
	}

})();
