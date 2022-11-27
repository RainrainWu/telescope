package telescope

import (
	"errors"
	"fmt"
	"io/ioutil"
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

type OutdatedScope int

const (
	MAJOR OutdatedScope = iota
	MINOR
	PATCH
	UP_TO_DATE
)

type IReportable interface {
	queryVersionsInformation()
	ReportOutdated(scope OutdatedScope)
}

type Atlas struct {
	name         string
	language     Language
	dependencies []IDependable
}

type PoetryLockPackage struct {
	Name     string `toml:"name"`
	Version  string `toml:"version"`
	Category string `toml:"category"`
}

type PoetryLock struct {
	Packages []PoetryLockPackage `toml:"package"`
}

func NewAtlas(filePath string) IReportable {

	var atlas IReportable

	fileBytes := parseDependenciesFile(filePath)

	switch filePath {
	case "go.mod":
		atlas = buildAtlasGoMod(fileBytes)
	case "poetry.lock":
		atlas = buildAtlasPoetryLock(fileBytes)
	default:
		panic(errors.New(fmt.Sprintf("unknown dep file type: %s", filePath)))
	}
	atlas.queryVersionsInformation()
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

func buildAtlasGoMod(fileBytes []byte) IReportable {

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
		if require.Indirect {
			continue
		}
		atlas.appendDependency(NewDependency(require.Mod.Path, require.Mod.Version))
	}
	return &atlas
}

func buildAtlasPoetryLock(fileBytes []byte) IReportable {

	var poetryLock PoetryLock
	toml.Unmarshal(fileBytes, &poetryLock)

	atlas := Atlas{
		name:         "",
		language:     PYTHON,
		dependencies: []IDependable{},
	}
	for _, pkg := range poetryLock.Packages {
		if pkg.Category != "main" {
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

func (a *Atlas) ReportOutdated(scope OutdatedScope) {

	outdatedMap := map[OutdatedScope][]IDependable{
		MAJOR: []IDependable{},
		MINOR: []IDependable{},
		PATCH: []IDependable{},
		UP_TO_DATE: []IDependable{},
	}

	for _, dep := range a.dependencies {
		depOutdatedScope := dep.(*Dependency).GetOutdatedScope()
		outdatedMap[depOutdatedScope] = append(outdatedMap[depOutdatedScope], dep)
	}
	for scp, deps := range outdatedMap {
		if scp > scope {
			continue
		}
		for _, dep := range deps {
			fmt.Println(dep.(*Dependency).Name, dep.(*Dependency).VersionCurrent, dep.(*Dependency).VersionLatest)
		}
	}
}
