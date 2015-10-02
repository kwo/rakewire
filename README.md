# Rakewire

Welcome to the Rakewire source code, developer!

This source code is divided into two projects: the backend, written in go, and the UI, an SPA, written using React.

## GO

### Prepare Environment

The directory $PROJECT_ROOT/go acts as the GOPATH and contains the bin, pkg and src directories.
The project dependencies are located at vendor/src according to the gb convention.

Therefore, set the following environment varibles:

	export GOPATH=$PROJECT_ROOT/go:$PROJECT_ROOT/go/vendor
	export GOBIN=$PROJECT_ROOT/go/bin

Additionally add GOBIN and web/npm_modules/.bin to the PATH

	export PATH=$GOBIN:$PROJECT_ROOT/web/npm_modules/.bin:$PATH

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

	cd $PROJECT_ROOT/web
	npm install

### Build

To compile a UI fit for production.

	webpack

### Run

Run webpack in watch mode to keep files up-to-date.

	cd $PROJECT_ROOT
	./build.sh web
