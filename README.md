# Rakewire

## Building

	git clone --recursive https://code.kfabrik.de:3333/rakewire/rakewire
	go test $(go list ./... | grep -v /vendor/)
	#go generate $(go list ./... | grep -v /vendor/)

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 LDFLAGS="-X main.Version=1.12.2 -X main.BuildTime=`date -u +%FT%TZ` -X main.BuildHash=`git rev-parse HEAD`"; go install -tags netgo -ldflags "$LDFLAGS" rakewire.go

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


## Test Service

	curl -D - -u karl@ostendorf.com:abcdefg -H "Content-Type: application/json" -d '{}' https://rw.kfabrik.de:8888/api/status -X POST
	/usr/local/Cellar/curl/7.48.0/bin/curl -D - -u karl@ostendorf.com:abcdefg -H "Content-Type: application/json" -d '{}' https://rw.kfabrik.de:8888/api/status -X POST
	curl -D - -H "Content-Type: application/json" -d '{}' https://rw.kfabrik.de:8888/api/status -X POST
	curl -D - -u karl@ostendorf.com:abcdefg https://rw.kfabrik.de:8888/api/status -X POST

	curl -D - -u ko:abcdefg https://rw.kfabrik.de:8888/subscriptions.opml

	/usr/local/Cellar/curl/7.48.0/bin/curl -D - -u karl@ostendorf.com:abcdefg -H "Content-Type: application/json" -d '{}' https://rakewire.kfabrik.de/api/status -X POST


## OPML

curl -u karl@ostendorf.com:abcdefg https://${RAKEWIRE_INSTANCE}/subscriptions.opml > rakewire.opml
curl -u karl@ostendorf.com:abcdefg -X PUT --data-binary @rakewire.opml https://${RAKEWIRE_INSTANCE}/subscriptions.opml

### obsolete
curl -u karl@ostendorf.com:abcdefg https://rakewire.kfabrik.de/api/rakewire.opml > rakewire.opml
curl -u karl@ostendorf.com:abcdefg -X PUT --data-binary @rakewire.opml https://rakewire.kfabrik.de/api/rakewire.opml
