# Rakewire

## Building

	git clone https://code.kfabrik.de:3333/rakewire/rakewire
	cd rakewire
	git submodule update --init
	go test $(go list ./... | grep -v /vendor/)

	go install ./tools/gokv/gokv.go
	go install ./tools/esc
	go generate $(go list ./... | grep -v /vendor/)
	go test $(go list ./... | grep -v /vendor/)

	go install rakewire.go

	CGO_ENABLED=0; LDFLAGS="-X main.Version=1.12.0-beta -X main.BuildTime=`date -u +%FT%TZ` -X main.BuildHash=`git rev-parse HEAD`"; go install -tags netgo -ldflags "$LDFLAGS" rakewire.go

## Dependencies

Dependencies are managed using git submodules via the vendetta tool.

One tool, esc, must be installed manually to the tools directory as follows

	cd tools
	git submodule add https://github.com/mjibson/esc

additionally, vendetta does not install dependencies of test files by default so they must be installed manually as well

 	cd vendor/github.com
	mkdir antonholmquist
	cd antonholmquist
	git submodule add https://github.com/antonholmquist/jason


## Test Echo service

		curl -X POST https://rw.kfabrik.de:4444/api/status -H "Content-Type: application/json" -d '{}'


## OPML

### Local

curl -D - -u karl@ostendorf.com:abcdefg http://localhost:8888/api/rakewire.opml > rakewire.opml
curl -D - -u karl@ostendorf.com:abcdefg -X PUT --data-binary @rakewire.opml http://localhost:8888/api/rakewire.opml

### Production

curl -u karl@ostendorf.com:abcdefg https://rakewire.kfabrik.de/api/rakewire.opml > rakewire.opml
curl -u karl@ostendorf.com:abcdefg -X PUT --data-binary @rakewire.opml https://rakewire.kfabrik.de/api/rakewire.opml
