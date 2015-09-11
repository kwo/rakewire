#! /bin/bash

CMD=$1

case $CMD in

	"clean")
		git clean -fdx
		;;

	"build")
		go build ./src/rakewire/rakewire.go
	  ;;

	"depgraph")
		godepgraph -s -horizontal rakewire | dot -Tsvg -o depgraph.svg && open depgraph.svg
	  ;;

	"install")
		go install ./src/rakewire/rakewire.go
	  ;;

	"run")
		go run ./src/rakewire/rakewire.go
	  ;;

	"test")
		go test ./src/...
	  ;;

	"update")
		gb vendor update --all
	  ;;

  *)
	  echo "unknown command: $CMD"
		echo "Usage `basename $0`: clean | build | depgraph | install | run | test | update"
		;;

esac
