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
	skipUnknown   bool
)

func init() {

	flag.StringVar(&filePath, "f", "go.mod", "dependencies file path")
	flag.StringVar(&outdatedScope, "s", "major", "desired outdated scope")
	flag.BoolVar(&skipUnknown, "skip-unknown", false, "skip dependencies with unknown versions")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: telescope [-f file_path] [-s outdated_scope] [--skip-unknown]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	atlas := telescope.NewAtlas(filePath)
	atlas.ReportOutdated(
		telescope.OutdatedScopeStrToEnum(outdatedScope),
		skipUnknown,
	)
}
