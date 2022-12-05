package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"telescope/telescope"
)

type IgnoredExpressions map[string]bool

func (i *IgnoredExpressions) String() string {

	return fmt.Sprintln(i.ToSlice())
}

func (i *IgnoredExpressions) Set(value string) error {

	(*i)[value] = true
	return nil
}

func (i *IgnoredExpressions) ToSlice() []string {

	var ignored []string
	for dep := range *i {
		ignored = append(ignored, dep)
	}
	return ignored
}

type CriticalExpressions map[string]telescope.OutdatedScope

func (c *CriticalExpressions) String() string {

	var buffer string
	criticalMap := c.ToScopeMap()
	for _, scp := range telescope.OutdatedScopeSeries {

		buffer += fmt.Sprintf("[ %s ]\n%s\n", scp, criticalMap[scp])
	}
	return buffer
}

func (c *CriticalExpressions) Set(value string) error {

	desiredScopeStr, expression, found := strings.Cut(value, ":")
	expressionMap := map[string]telescope.OutdatedScope(*c)
	if !found {
		desiredScopeStr, expression = telescope.MAJOR.String(), desiredScopeStr
	}

	desiredScope := telescope.OutdatedScopeStrToEnum(desiredScopeStr)
	if registeredScope, ok := expressionMap[expression]; ok {
		expressionMap[expression] = telescope.GetTopScope(
			[]telescope.OutdatedScope{registeredScope, desiredScope},
		)
	} else {
		expressionMap[expression] = desiredScope
	}
	return nil
}

func (c *CriticalExpressions) ToScopeMap() map[telescope.OutdatedScope][]string {

	var criticalMap map[telescope.OutdatedScope][]string = map[telescope.OutdatedScope][]string{}
	for dep, scp := range *c {

		criticalMap[scp] = append(criticalMap[scp], dep)
	}
	return criticalMap
}

var (
	filePath            string
	outdatedScope       string
	skipUnknown         bool
	strictSemVer        bool
	ignoredExpressions  IgnoredExpressions  = make(map[string]bool)
	criticalExpressions CriticalExpressions = make(map[string]telescope.OutdatedScope)
)

func init() {

	flag.StringVar(&filePath, "f", "go.mod", "dependencies file path")
	flag.StringVar(&outdatedScope, "s", "major", "desired outdated scope")
	flag.BoolVar(&skipUnknown, "skip-unknown", false, "skip dependencies with unknown versions")
	flag.BoolVar(&strictSemVer, "strict-semver", false, "parse dependencies file with strict SemVer format")
	flag.Var(&ignoredExpressions, "i", "ignore specific dependency")
	flag.Var(&criticalExpressions, "c", "highlight critical dependency")
	flag.Usage = usage
}

func usage() {

	fmt.Fprintf(os.Stderr, "Usage: telescope [-f file_path] [-s outdated_scope] [-i ignored_dependency] [-c critical_dependency] [--skip-unknown] [--strict-semver]\n")
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	atlas := telescope.NewAtlas(filePath, strictSemVer, ignoredExpressions.ToSlice(), criticalExpressions.ToScopeMap())
	criticalFound := atlas.ReportOutdated(
		telescope.OutdatedScopeStrToEnum(outdatedScope),
		skipUnknown,
	)
	if criticalFound {
		os.Exit(1)
	}
}
