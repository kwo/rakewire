# Rakewire

## Building

	git clone --recursive https://github.com/kwo/rakewire
	go test $(go list ./... | grep -v /vendor/)
	#go generate $(go list ./... | grep -v /vendor/)

	linux:
	export GOOS=linux
	export GOARCH=amd64
	export CGO_ENABLED=0
	export LDFLAGS="-X main.Version=$(cat VERSION) -X main.BuildTime=`date -u +%FT%TZ` -X main.BuildHash=`git rev-parse HEAD`"
	go install -tags netgo -ldflags "$LDFLAGS" rakewire.go

	macOS:
	export CGO_ENABLED=0
	export LDFLAGS="-X main.Version=$(cat VERSION) -X main.BuildTime=`date -u +%FT%TZ` -X main.BuildHash=`git rev-parse HEAD`"
	go install -tags netgo -ldflags "$LDFLAGS" rakewire.go

## Dependencies

Dependencies are managed using git submodules via the [vendetta](https://github.com/dpw/vendetta) tool.
