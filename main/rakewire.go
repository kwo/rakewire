package main

import (
	"fmt"
	"os"
	"path"
	"rakewire.com/fetch"
	"strings"
)

func main() {

	var param string
	if len(os.Args[1:]) == 1 {
		param = strings.TrimSpace(os.Args[1])
	}
	if len(param) == 0 {
		fmt.Printf("Usage: %s <feed-file>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	fetch.Fetch(param)

}
