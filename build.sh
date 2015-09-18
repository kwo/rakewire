#! /bin/bash

CMD=$1

case $CMD in

	"clean")
		git clean -fdx
		;;

	"build")
		#cd web
		#gulp build
		#gulp buildmode
		go build ./src/rakewire/rakewire.go
		rice append --exec rakewire -i ./src/rakewire/httpd
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
