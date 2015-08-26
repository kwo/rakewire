(function() {

	'use strict';

	const gulp = require('gulp');
	const gutil = require('gulp/node_modules/gulp-util');
	const plumber = require('gulp-plumber');
	const paths = makePaths();

	gulp.task('default', ['monitor']);
	gulp.task('lint', lint);
	gulp.task('clean', clean);
	gulp.task('build', ['clean', 'lint'], build);
	gulp.task('monitor', ['build'], monitor);
	gulp.task('deploy', ['build'], deploy);

	function makePaths() {

		let base = 'app';
		let src = base + '/src';
		let dst = base + '/lib';

		let p = {
			dst: {
				base: dst,
				all:  dst + '/**/*'
			},
			src: {
				base: src,
				all:  src + '/**/*',
				js:   src + '/**/*.js',
				static: [src + '/**/*', '!' + src + '/**/*.js']
			},
			vendor: {
				all: base + '/modules/**/*',
				base: base + '/modules',
				src: 'node_modules'
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

	function build() {
		return Promise.all([resources(), transpile()]);
	}

	function bundle() {
		log('bundling certain JSPM packages...');
		const jspm = require('jspm');
		const path = require('path');
		return jspm.bundle(
			'react + react-router + react-bootstrap + react-router-bootstrap + fetch',
			path.join(paths.dst.base, 'bundle.js'),
			{ inject: true, mangle: true, minify: true, sourceMaps: false }
		);
	}

	function resources() {
		log('copying static resources...');
		return new Promise(function(resolve, reject) {
			gulp.src(paths.src.static)
				.pipe(gulp.dest(paths.dst.base))
				.on('end', resolve)
				.on('error', reject);
		});
	}

	function transpile() {
		log('transpiling to ES5...');
		const babel = require('gulp-babel');
		return new Promise(function(resolve, reject) {
			gulp.src(paths.src.js)
				.pipe(plumber())
				.pipe(babel())
				.pipe(gulp.dest(paths.dst.base))
				.on('end', resolve)
				.on('error', reject);
		});
	}

	function monitor(done) {

		done = done || {}; // never call done

		log('monitoring...');

		const babel = require('gulp-babel');
		const del = require('del');
		const ignore = require('gulp-ignore');
		const path = require('path');

		function cleanFile(filename) {
			let dstFile = path.join(paths.dst.base, filename);
			return new Promise(function(resolve, reject) {
				del(dstFile, function(err) {
					if (err) return reject(err);
					resolve();
				});
			});
		}

		function buildFile(filename) {
			if (path.extname(filename) === '.js') {
				return transpileFile(filename);
			} else {
				return copyFile(filename);
			}
		}

		function transpileFile(filename) {
			return new Promise(function(resolve, reject) {
				gulp.src(paths.src.all)
					.pipe(plumber())
					.pipe(ignore.include(function(file) {
						return (file.relative === filename);
					}))
					.pipe(babel())
					.pipe(gulp.dest(paths.dst.base))
					.on('end', resolve)
					.on('error', reject);
			});
		}

		function copyFile(filename) {
			return new Promise(function(resolve, reject) {
				gulp.src(paths.src.all)
					.pipe(ignore.include(function(file) {
						return (file.relative === filename);
					}))
					.pipe(gulp.dest(paths.dst.base))
					.on('end', resolve)
					.on('error', reject);
			});
		}

		let watcher = gulp.watch(paths.src.all);
		watcher.on('change', function watch(event) {
			let filename = path.relative(paths.src.base, event.path);
			log(event.type, filename);
			if (['added', 'changed', 'renamed'].indexOf(event.type) !== -1) {
				buildFile(filename);
			} else if (event.type === 'deleted') {
				cleanFile(filename);
			} else {
				log('Unknown event type:', event.type);
			}
		});

	}

	function deploy() {
		log('not implemented');
	}

})();
