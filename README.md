# Rakewire

Welcome to the Rakewire source code, developer!

This source code is divided into two projects: the backend, written in go, and the UI, an SPA, written using React.

## GO

### Prepare Environment

The project root directory ($PROJECT_ROOT) acts as the GOPATH and contains the bin, pkg and src directories.
The project dependencies are located at src/rakewire/vendor according to the GO15VENDOREXPERIMENT=1 convention.

Therefore, set the following environment varibles:

	export GOPATH=$PROJECT_ROOT
	export GOBIN=$GOPATH/bin
	export GO15VENDOREXPERIMENT=1

### Build

	cd $PROJECT_ROOT
	go build src/rakewire/rakewire.go

### Run

	cd $PROJECT_ROOT
	go run src/rakewire/rakewire.go

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
