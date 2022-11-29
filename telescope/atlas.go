package telescope

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	queryVersionsInformation()
	buildOutdatedMap()
	ReportOutdated(scope OutdatedScope, skipUnknown bool)
}

type Atlas struct {
	name         string
	language     Language
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

func NewAtlas(filePath string, ignoredDeps map[string]bool) IReportable {

	var atlas IReportable

	fileBytes := parseDependenciesFile(filePath)
	splitPath := strings.Split(filePath, "/")
	fileName := splitPath[len(splitPath)-1]

	switch fileName {
	case "go.mod":
		atlas = buildAtlasGoMod(fileBytes, ignoredDeps)
	case "poetry.lock":
		atlas = buildAtlasPoetryLock(fileBytes, ignoredDeps)
	default:
		panic(errors.New(fmt.Sprintf("unknown dep file: %s", filePath)))
	}
	atlas.queryVersionsInformation()
	atlas.buildOutdatedMap()
	return atlas
}

func parseDependenciesFile(filePath string) []byte {

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Fatal(err.Error())
		panic(errors.New(fmt.Sprintf("failed to read dep file %s", filePath)))
	}

	return fileBytes
}

func buildAtlasGoMod(fileBytes []byte, ignoredDeps map[string]bool) IReportable {

	modObject, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	atlas := Atlas{
		name:         modObject.Module.Mod.Path,
		language:     GO,
		dependencies: []IDependable{},
	}
	for _, require := range modObject.Require {
		if _, ok := ignoredDeps[require.Mod.Path]; ok {
			continue
		}
		atlas.appendDependency(NewDependency(require.Mod.Path, require.Mod.Version))
	}
	return &atlas
}

func buildAtlasPoetryLock(fileBytes []byte, ignoredDeps map[string]bool) IReportable {

	var poetryLock PoetryLock
	toml.Unmarshal(fileBytes, &poetryLock)

	atlas := Atlas{
		name:         "",
		language:     PYTHON,
		dependencies: []IDependable{},
		outdatedMap:  map[OutdatedScope][]IDependable{},
	}
	for _, pkg := range poetryLock.Packages {
		if _, ok := ignoredDeps[pkg.Name]; ok {
			continue
		}
		atlas.appendDependency(NewDependency(pkg.Name, pkg.Version))
	}
	return &atlas
}

func (a *Atlas) appendDependency(dep IDependable) {

	a.dependencies = append(a.dependencies, dep)
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

func (a *Atlas) ReportOutdated(desiredScope OutdatedScope, skipUnknown bool) {

	for _, scp := range [3]OutdatedScope{MAJOR, MINOR, PATCH} {
		if scp > desiredScope {
			break
		}
		a.reportByScope(scp)
	}
	if !skipUnknown {
		a.reportUnknownDependencies()
	}
}

func buildReportItem(dep IDependable) string {

	if dep.(*Dependency).VersionCurrent == nil || dep.(*Dependency).VersionLatest == nil {
		return fmt.Sprintf(
			"%-40s %-20s",
			dep.(*Dependency).Name,
			dep.(*Dependency).VersionCurrentLiteral,
		)
	}
	return fmt.Sprintf(
		"%-40s %-20s %-20s",
		dep.(*Dependency).Name,
		dep.(*Dependency).VersionCurrent,
		dep.(*Dependency).VersionLatest,
	)
}

func (a *Atlas) reportByScope(scope OutdatedScope) {

	fmt.Printf(
		"\n[%s Version Outdated]%s\n",
		scope.String(),
		strings.Repeat("=", 20),
	)
	if len(a.outdatedMap[scope]) == 0 {
		fmt.Printf("no outdated dependencies")
	}
	for _, dep := range a.outdatedMap[scope] {
		fmt.Println(buildReportItem(dep))
	}
	fmt.Print("\n")
}

func (a *Atlas) reportUnknownDependencies() {

	if len(a.outdatedMap[UNKNOWN]) == 0 {
		return
	}
	fmt.Printf("\n[UNKNOWN dependencies]%s\n", strings.Repeat("=", 20))
	for _, dep := range a.outdatedMap[UNKNOWN] {
		fmt.Println(buildReportItem(dep))
	}
}
