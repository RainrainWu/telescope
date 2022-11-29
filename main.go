package main

import (
	"flag"
	"fmt"
	"os"
	"telescope/telescope"
)

type IgnoredDependencies map[string]bool

func (i *IgnoredDependencies) String() string {

	var ignored []string
	for dep := range *i {
		ignored = append(ignored, dep)
	}
	return fmt.Sprintln(ignored)
}

func (i *IgnoredDependencies) Set(value string) error {

	map[string]bool(*i)[value] = true
	return nil
}

var (
	filePath            string
	outdatedScope       string
	skipUnknown         bool
	strictSemVer        bool
	ignoredDependencies IgnoredDependencies = make(map[string]bool)
)

func init() {

	flag.StringVar(&filePath, "f", "go.mod", "dependencies file path")
	flag.StringVar(&outdatedScope, "s", "major", "desired outdated scope")
	flag.BoolVar(&skipUnknown, "skip-unknown", false, "skip dependencies with unknown versions")
	flag.BoolVar(&strictSemVer, "strict-semver", false, "parse dependencies file with strict SemVer format")
	flag.Var(&ignoredDependencies, "i", "ignore specific dependency")
	flag.Usage = usage
}

func usage() {

	fmt.Fprintf(os.Stderr, "Usage: telescope [-f file_path] [-s outdated_scope] [-i ignored_dependency] [--skip-unknown]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	atlas := telescope.NewAtlas(filePath, strictSemVer, ignoredDependencies)
	atlas.ReportOutdated(
		telescope.OutdatedScopeStrToEnum(outdatedScope),
		skipUnknown,
	)
}
