# Rakewire

Welcome to the Rakewire source code, developer!

This source code is divided into two projects: the backend, written in go, and the UI, an SPA, written using React.

## GO

### Prepare Environment

The project root directory ($PROJECT_ROOT) acts as the GOPATH and contains the bin, pkg and src directories.
The project dependencies are located at vendor/src according to the gb convention.

Therefore, set the following environment varibles:

	export GOPATH=$PROJECT_ROOT:$PROJECT_ROOT/vendor
	export GOBIN=$PROJECT_ROOT/bin

### Build

	cd $PROJECT_ROOT
	./build.sh build

### Run

	cd $PROJECT_ROOT
	./build.sh run

### Top-Level Dependencies

 - [github.com/GeertJohan/go.rice](https://github.com/GeertJohan/go.rice)
 - [github.com/boltdb/bolt](https://github.com/boltdb/bolt)
 - [github.com/gorilla/handlers](https://github.com/gorilla/handlers)
 - [github.com/gorilla/mux](https://github.com/gorilla/mux)
 - [github.com/pborman/uuid](https://github.com/pborman/uuid)
 - [github.com/rogpeppe/go-charset/charset](https://github.com/rogpeppe/go-charset/charset)
 - [github.com/rogpeppe/go-charset/data](https://github.com/rogpeppe/go-charset/data)
 - [github.com/stretchr/testify/assert](https://github.com/stretchr/testify)
 - [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)

## Web

### Prepare Environment

Install the node modules and jspm packages

	npm install

### JSPM

Sometimes dependencies will upgrade react to a still-in-beta version. Therefore is must be set back as follows

	jspm resolve --only npm:react@0.13.3
	jspm clean

### Run

Simple accessing the UI at the test URL (localhost:4444) will cause JSPM to compile the application with Babel and JSX.

	cd $PROJECT_ROOT
	./build.sh run

### Build

To compile a UI fit for production.

	gulp build

Use `gulp buildmode` to access the production-mode version. Use `gulp devmode` to switch back.
