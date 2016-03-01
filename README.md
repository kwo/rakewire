# Rakewire

Welcome to the Rakewire source code, developer!

This source code is divided into two projects: the backend, written in go, and the UI, an SPA, written using React.

## Dependencies

Dependencies are managed using git submodules via the vendetta tool.

One tool, esc, must be installed manually to the tools directory as follows

	cd tools
	git submodule add https://github.com/mjibson/esc

additionally, vendetta does not install dependencies of test file files by default so they must be installed manually as well

 	cd vendor/github.com
	mkdir antonholmquist
	cd antonholmquist
	git submodule add https://github.com/antonholmquist/jason


## OPML

### Local

curl -D - -u karl@ostendorf.com:abcdefg http://localhost:8888/api/rakewire.opml > rakewire-dev.opml
curl -D - -u karl@ostendorf.com:abcdefg -X PUT --data-binary @rakewire.opml http://localhost:8888/api/rakewire.opml
curl -D - -u karl@ostendorf.com:abcdefg -D - -X POST http://localhost:8888/api/cleanup

### Production

curl -u karl@ostendorf.com:abcdefg https://rakewire.kfabrik.de/api/rakewire.opml > rakewire.opml
curl -u karl@ostendorf.com:abcdefg -X PUT --data-binary @rakewire.opml https://rakewire.kfabrik.de/api/rakewire.opml?replace=true
curl -u karl@ostendorf.com:abcdefg -X POST https://rakewire.kfabrik.de/api/cleanup
