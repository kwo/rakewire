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

	github.com/GeertJohan/go.rice
	github.com/boltdb/bolt
	github.com/codegangsta/negroni
	github.com/gorilla/mux
	github.com/pborman/uuid
	github.com/phyber/negroni-gzip/gzip
	github.com/rogpeppe/go-charset/charset
	github.com/rogpeppe/go-charset/data
	github.com/stretchr/testify/assert
	gopkg.in/yaml.v2


## UI

### Prepare Environment

Install the node modules and jspm packages

	npm install


### Run

Simple accessing the UI at the test URL (localhost:4444) will cause JSPM to compile the application with Babel and JSX.

	cd $PROJECT_ROOT
	./build.sh run

### Build

To compile a UI fit for production.

	gulp build

Use `gulp buildmode` to access the production-mode version. Use `gulp devmode` to switch back.
