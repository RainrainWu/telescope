package telescope

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"
)

type Language int

const (
	GO Language = iota
	PYTHON
)

func (l Language) String() string {
	return [...]string{"GO", "PYTHON"}[l]
}

type IReportable interface {
	ReportOutdated(scope OutdatedScope, skipUnknown bool) bool
}

type Atlas struct {
	name         string
	language     Language
	criticalMap  map[string]OutdatedScope
	dependencies []IDependable
	outdatedMap  map[OutdatedScope][]IDependable
}

type PoetryLockPackage struct {
	Name     string `toml:"name"`
	Version  string `toml:"version"`
	Category string `toml:"category"`
}

type PoetryLock struct {
	Packages []PoetryLockPackage `toml:"package"`
}

func NewAtlas(
	filePath string,
	strictSemVer bool,
	ignoredDeps map[string]bool,
	criticalDeps map[string]OutdatedScope,
) IReportable {

	var atlas IReportable

	fileBytes := parseDependenciesFile(filePath)
	splitPath := strings.Split(filePath, "/")
	fileName := splitPath[len(splitPath)-1]

	switch fileName {
	case "go.mod":
		atlas = buildAtlasGoMod(fileBytes, strictSemVer, ignoredDeps, criticalDeps)
	case "poetry.lock":
		atlas = buildAtlasPoetryLock(fileBytes, strictSemVer, ignoredDeps, criticalDeps)
	default:
		panic(fmt.Errorf("unknown dep file: %s", filePath))
	}
	atlas.(*Atlas).sortLexicographically()
	atlas.(*Atlas).queryVersionsInformation()
	atlas.(*Atlas).buildOutdatedMap()
	return atlas
}

func parseDependenciesFile(filePath string) []byte {

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Fatal(err.Error())
		panic(fmt.Errorf("failed to read dep file %s", filePath))
	}

	return fileBytes
}

func buildAtlasGoMod(fileBytes []byte, strictSemVer bool, ignoredDeps map[string]bool, criticalDeps map[string]OutdatedScope) IReportable {

	modObject, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	atlas := Atlas{
		name:         modObject.Module.Mod.Path,
		language:     GO,
		criticalMap:  criticalDeps,
		dependencies: []IDependable{},
	}
	for _, require := range modObject.Require {
		if _, ok := ignoredDeps[require.Mod.Path]; ok {
			continue
		}
		atlas.appendDependency(
			NewDependency(require.Mod.Path, require.Mod.Version, strictSemVer),
		)
	}
	return &atlas
}

func buildAtlasPoetryLock(fileBytes []byte, strictSemVer bool, ignoredDeps map[string]bool, criticalDeps map[string]OutdatedScope) IReportable {

	var poetryLock PoetryLock
	err := toml.Unmarshal(fileBytes, &poetryLock)
	if err != nil {
		panic(err)
	}

	atlas := Atlas{
		name:         "",
		language:     PYTHON,
		dependencies: []IDependable{},
		criticalMap:  criticalDeps,
		outdatedMap:  map[OutdatedScope][]IDependable{},
	}
	for _, pkg := range poetryLock.Packages {
		if _, ok := ignoredDeps[pkg.Name]; ok {
			continue
		}
		atlas.appendDependency(
			NewDependency(pkg.Name, pkg.Version, strictSemVer),
		)
	}
	return &atlas
}

func (a *Atlas) appendDependency(dep IDependable) {

	a.dependencies = append(a.dependencies, dep)
}

func (a *Atlas) sortLexicographically() {

	sort.SliceStable(
		a.dependencies,
		func(i, j int) bool {
			return strings.Compare(
				a.dependencies[i].(*Dependency).Name,
				a.dependencies[j].(*Dependency).Name,
			) == -1
		},
	)
}

func (a *Atlas) queryVersionsInformation() {

	queryWaitGroup := new(sync.WaitGroup)

	queryWaitGroup.Add(len(a.dependencies))
	for _, dep := range a.dependencies {
		go dep.QueryReleaseVersions(a.language, queryWaitGroup)
	}
	queryWaitGroup.Wait()
}

func (a *Atlas) buildOutdatedMap() {

	outdatedMap := map[OutdatedScope][]IDependable{
		MAJOR:      {},
		MINOR:      {},
		PATCH:      {},
		UP_TO_DATE: {},
		UNKNOWN:    {},
	}
	for _, dep := range a.dependencies {
		if dep.(*Dependency).VersionCurrent == nil {
			outdatedMap[UNKNOWN] = append(outdatedMap[UNKNOWN], dep)
			continue
		}
		depOutdatedScope := dep.(*Dependency).GetOutdatedScope()
		outdatedMap[depOutdatedScope] = append(outdatedMap[depOutdatedScope], dep)
	}
	a.outdatedMap = outdatedMap
}

func (a *Atlas) ReportOutdated(desiredScope OutdatedScope, skipUnknown bool) bool {

	var criticalFound bool

	for _, scp := range [3]OutdatedScope{MAJOR, MINOR, PATCH} {
		if scp > desiredScope {
			break
		}
		color := MapScopeColor[scp]
		criticalFound = a.reportByScope(scp, color) || criticalFound
	}
	if !skipUnknown {
		a.reportUnknownDependencies()
	}

	return criticalFound
}

func buildReportItem(dep IDependable) string {

	if dep.(*Dependency).VersionCurrent == nil || dep.(*Dependency).VersionLatest == nil {
		return fmt.Sprintf(
			"%-50s %-20s",
			dep.(*Dependency).Name,
			dep.(*Dependency).VersionCurrentLiteral,
		)
	}
	return fmt.Sprintf(
		"%-50s %-20s %-20s",
		dep.(*Dependency).Name,
		dep.(*Dependency).VersionCurrent,
		dep.(*Dependency).VersionLatest,
	)
}

func (a *Atlas) reportByScope(scope OutdatedScope, color int) bool {

	var criticalFound bool

	fmt.Printf(
		"\033[%dm\n[ %d %s Version Outdated ]%s\n\n",
		color,
		len(a.outdatedMap[scope]),
		scope.String(),
		strings.Repeat("=", 40),
	)
	if len(a.outdatedMap[scope]) == 0 {
		fmt.Printf("no outdated dependencies")
	}

	for _, dep := range a.outdatedMap[scope] {

		if scp, ok := a.criticalMap[dep.(*Dependency).Name]; ok && scp >= scope {
			criticalFound = criticalFound || true
			fmt.Printf("* %s\n", buildReportItem(dep))
		} else {
			fmt.Printf("  %s\n", buildReportItem(dep))
		}
	}
	fmt.Print("\n\033[0m")

	return criticalFound
}

func (a *Atlas) reportUnknownDependencies() {

	if len(a.outdatedMap[UNKNOWN]) == 0 {
		return
	}
	fmt.Printf(
		"\n[ %d UNKNOWN dependencies ]%s\n\n",
		len(a.outdatedMap[UNKNOWN]),
		strings.Repeat("=", 40),
	)
	for _, dep := range a.outdatedMap[UNKNOWN] {
		fmt.Printf("  %s\n", buildReportItem(dep))
	}
}
