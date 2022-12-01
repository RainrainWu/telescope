package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

type CriticalDependencies map[string]telescope.OutdatedScope

func (c *CriticalDependencies) String() string {

	var buffer string
	for dep, scp := range map[string]telescope.OutdatedScope(*c) {

		buffer += fmt.Sprintf("[ %s ] %s\n", scp, dep)
	}
	return buffer
}

func (c *CriticalDependencies) Set(value string) error {

	var depName string
	var desiredScope telescope.OutdatedScope

	result := strings.Split(value, ":")
	c_deps := map[string]telescope.OutdatedScope(*c)
	if len(result) == 2 {
		depName, desiredScope = result[1], telescope.OutdatedScopeStrToEnum(result[0])
	} else if len(result) == 1 {
		depName, desiredScope = result[0], telescope.MAJOR
	} else {
		panic(fmt.Errorf("invalid expression: %s", value))
	}

	if registeredScope, ok := c_deps[depName]; ok {
		c_deps[depName] = telescope.GetTopScope(
			[]telescope.OutdatedScope{registeredScope, desiredScope},
		)
	} else {
		c_deps[depName] = desiredScope
	}
	return nil
}

var (
	filePath             string
	outdatedScope        string
	skipUnknown          bool
	strictSemVer         bool
	ignoredDependencies  IgnoredDependencies  = make(map[string]bool)
	criticalDependencies CriticalDependencies = map[string]telescope.OutdatedScope{"*": telescope.MAJOR}
)

func init() {

	flag.StringVar(&filePath, "f", "go.mod", "dependencies file path")
	flag.StringVar(&outdatedScope, "s", "major", "desired outdated scope")
	flag.BoolVar(&skipUnknown, "skip-unknown", false, "skip dependencies with unknown versions")
	flag.BoolVar(&strictSemVer, "strict-semver", false, "parse dependencies file with strict SemVer format")
	flag.Var(&ignoredDependencies, "i", "ignore specific dependency")
	flag.Var(&criticalDependencies, "c", "highlight critical dependency")
	flag.Usage = usage
}

func usage() {

	fmt.Fprintf(os.Stderr, "Usage: telescope [-f file_path] [-s outdated_scope] [-i ignored_dependency] [-c critical_dependency] [--skip-unknown] [--strict-semver]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	atlas := telescope.NewAtlas(filePath, strictSemVer, ignoredDependencies, criticalDependencies)
	criticalFound := atlas.ReportOutdated(
		telescope.OutdatedScopeStrToEnum(outdatedScope),
		skipUnknown,
	)
	if criticalFound {
		os.Exit(1)
	}
}
