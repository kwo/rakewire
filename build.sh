#! /bin/bash

CMD=$1

case $CMD in

	"clean")
		git clean -fdx
		;;

	"build")

		echo "building webapp"
		cd web
		webpack
		cd ..

		echo "building go application"
		cd go
		go build ./src/rakewire/rakewire.go

		echo "embedding webapp in go application"
		rm -f src/rakewire/httpd/public
		mv ../web/public src/rakewire/httpd
		rice append --exec rakewire -i ./src/rakewire/httpd
		mv src/rakewire/httpd/public ../web
		cd src/rakewire/httpd
		ln -s ../../../../web/public
		cd ../../../..

		echo "moving final app to go/bin"
		mv go/rakewire go/bin

	  ;;

	"depgraph")
		cd go
		godepgraph -s -horizontal rakewire | dot -Tsvg -o depgraph.svg && open depgraph.svg
	  ;;

	"run")
		cd go
		go run ./src/rakewire/rakewire.go
	  ;;

	"test")
		cd go
		go test ./src/...
	  ;;

	"update")
		cd go
		gb vendor update --all
	  ;;

	"web")
		cd web
		webpack --debug --watch --color
	  ;;

  *)
	  echo "unknown command: $CMD"
		echo "Usage `basename $0`: clean | build | depgraph | run | test | update | web"
		;;

esac
