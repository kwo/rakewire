#! /bin/bash

CMD=$1

case $CMD in

	"clean")
		git clean -fdx
		;;

	"build")

		cd web
		webpack
		cd ..

		go build ./src/rakewire/rakewire.go

		rm -f src/rakewire/httpd/public
		mv web/public src/rakewire/httpd
		rice append --exec rakewire -i ./src/rakewire/httpd
		mv src/rakewire/httpd/public web
		cd src/rakewire/httpd
		ln -s ../../../web/public
		cd ../../..

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

	"web")
		cd web
		webpack --debug --watch --color
	  ;;

  *)
	  echo "unknown command: $CMD"
		echo "Usage `basename $0`: clean | build | depgraph | install | run | test | update"
		;;

esac
