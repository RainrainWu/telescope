package main

import (
	"flag"
	"fmt"
	"os"
	"telescope/telescope"
)

var (
	filePath string
)

func init() {

	flag.StringVar(&filePath, "f", "go.mod", "dependencies file path")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: telescope [-f file_path]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	atlas := telescope.NewAtlas(filePath)
	atlas.ReportOutdated(telescope.PATCH)
}
