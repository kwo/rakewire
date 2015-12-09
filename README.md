# Rakewire

Welcome to the Rakewire source code, developer!

This source code is divided into two projects: the backend, written in go, and the UI, an SPA, written using React.

## GO

### Prepare Environment

The directory $PROJECT_ROOT/go acts as the GOPATH and contains the bin, pkg and src directories.
The project dependencies are located at vendor/src according to the gb convention.

Therefore, set the following environment varibles:

	export GOPATH=$PROJECT_ROOT/go/vendor:$PROJECT_ROOT/go
	export GOBIN=$PROJECT_ROOT/go/vendor/bin

Additionally add GOBIN and web/npm_modules/.bin to the PATH

	export PATH=$GOBIN:$PROJECT_ROOT/web/npm_modules/.bin:$PATH

Finally, be sure the following go executables are installed are in the PATH.
They are necessary for go generate commands.

	go get -u github.com/mjibson/esc


### Build

This will place the executable in $PROJECT_ROOT

	b.build

### Run

	g.run

### Top-Level Dependencies

 - [github.com/boltdb/bolt](https://github.com/boltdb/bolt)
 - [github.com/gorilla/handlers](https://github.com/gorilla/handlers)
 - [github.com/gorilla/mux](https://github.com/gorilla/mux)
 - [github.com/paulrosania/go-charset/charset](https://github.com/paulrosania/go-charset/charset)
 - [github.com/paulrosania/go-charset/data](https://github.com/paulrosania/go-charset/data)
 - [github.com/pborman/uuid](https://github.com/pborman/uuid)
 - [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)

## Web

### Prepare Environment

Install the node modules.

	cd $PROJECT_ROOT/web
	npm install

### Build

To compile a UI fit for production.

	w.build

### Run

Run webpack in watch mode to keep files up-to-date.

	w.run
