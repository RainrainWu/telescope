package main

import (
	"flag"
	"fmt"
	"os"
	"telescope/telescope"
)

var (
	filePath      string
	outdatedScope string
)

func init() {

	flag.StringVar(&filePath, "f", "go.mod", "dependencies file path")
	flag.StringVar(&outdatedScope, "s", "major", "desired outdated scope")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: telescope [-f file_path] [-s outdated_scope]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	atlas := telescope.NewAtlas(filePath)
	atlas.ReportOutdated(telescope.OutdatedScopeStrToEnum(outdatedScope))
}