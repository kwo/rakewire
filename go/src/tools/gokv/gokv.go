package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"tools/gokv/generator"
)

func main() {

	flag.Parse()
	extra := flag.Args()

	if len(extra) != 1 {
		usage()
		os.Exit(1)
	}

	inFile := extra[0]

	if path.Ext(inFile) != ".go" {
		fmt.Printf("Not a go file: %s\n", inFile)
		os.Exit(1)
	}

	outFile := path.Join(path.Dir(inFile), strings.TrimSuffix(path.Base(inFile), path.Ext(inFile))+"_kv.go")

	if err := generator.Generate(inFile, outFile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\t%s [options] <filename>\n\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}
