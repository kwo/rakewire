#! /bin/bash

CMD=$1

case $CMD in

	"run")
		go run ./src/rakewire/rakewire.go
	  ;;

	"test")
		go test $(go list ./src/rakewire/... | grep -v /vendor/)
	  ;;

	"build")
		go build ./src/rakewire/rakewire.go
	  ;;

	"install")
		go install ./src/rakewire/rakewire.go
	  ;;

	"depgraph")
		godepgraph -s rakewire | dot -Tsvg -o depgraph.svg
		open depgraph.svg
	  ;;

  *)
	  echo "unknown command: $CMD"
		echo "Usage `basename $0`: run | test | build | install | depgraph"
		;;

esac
