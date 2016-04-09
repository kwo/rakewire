package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const buildInfoTpl = `
package model

// This file was automatically generated with the buildinfo tool.
// Do not modify and do not place under version control.

// BuildInfo variables
var (
	BuildHash = "{{.Hash}}{{.Dirty}}"
	BuildTime = "{{.Time}}"
)
`

// CommitInfo encapsulates the data captured from the last git commit
type CommitInfo struct {
	Hash  string
	Time  string
	Dirty string
}

func main() {

	info, err := getLastGitCommitInfo()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	tpl, err := template.New("tpl").Parse(buildInfoTpl)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var out bytes.Buffer

	if err := tpl.Execute(&out, info); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	data, err := format.Source(out.Bytes())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	outFile, err := os.OpenFile(getOutputFilename(), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	defer outFile.Close()

	if _, err := outFile.Write(data); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

}

func getGitDirtyStatus() bool {

	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return true
	}

	return len(strings.TrimSpace(out.String())) > 0

}

func getLastGitCommitInfo() (*CommitInfo, error) {

	cmd := exec.Command("git", "log", "-1", "--format='%H %cI'")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	line := strings.Trim(out.String(), "'\n")
	fields := strings.Fields(line)
	dirty := ""
	if getGitDirtyStatus() {
		dirty = "*"
	}

	return &CommitInfo{Hash: fields[0], Time: fields[1], Dirty: dirty}, nil

}

func getOutputFilename() string {
	outputFilename := "buildinfo.go"
	args := os.Args
	if len(args) > 1 {
		outputFilename = args[1]
	}
	return outputFilename
}
