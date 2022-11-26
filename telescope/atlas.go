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

type Atlas struct {
	name         string
	dependencies []*Dependency
}

type PoetryLockPackage struct {
	Name     string `toml:"name"`
	Version  string `toml:"version"`
	Category string `toml:"category"`
}

type PoetryLock struct {
	Packages []PoetryLockPackage `toml:"package"`
}

func NewAtlas(filePath string) Atlas {

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	if filePath == "go.mod" {
		return ParseAtlasGoMod(fileBytes)
	} else if filePath == "poetry.lock" {
		return ParseAtlasPoetryLock(fileBytes)
	}
	panic(errors.New("unknown file type"))
}

func ParseAtlasGoMod(fileBytes []byte) Atlas {

	modObject, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	atlas := Atlas{
		name:         modObject.Module.Mod.Path,
		dependencies: []*Dependency{},
	}
	for _, require := range modObject.Require {
		if require.Indirect {
			continue
		}
		atlas.dependencies = append(
			atlas.dependencies,
			NewDependency(require.Mod.Path, require.Mod.Version, GO),
		)
	}
	return atlas
}

func ParseAtlasPoetryLock(fileBytes []byte) Atlas {

	var poetryLock PoetryLock
	toml.Unmarshal(fileBytes, &poetryLock)

	atlas := Atlas{
		name:         "",
		dependencies: []*Dependency{},
	}
	for _, pkg := range poetryLock.Packages {
		if pkg.Category != "main" {
			continue
		}
		fmt.Println(pkg.Name, pkg.Version)
		atlas.dependencies = append(
			atlas.dependencies,
			NewDependency(pkg.Name, pkg.Version, PYTHON),
		)
	}
	return atlas
}

func (a *Atlas) Query() {

	queryWaitGroup := new(sync.WaitGroup)

	queryWaitGroup.Add(len(a.dependencies))
	for _, dep := range a.dependencies {
		go dep.QueryVersionsPython(queryWaitGroup)
	}
	queryWaitGroup.Wait()
}

func (a *Atlas) Report() {

	for _, dep := range a.dependencies {
		fmt.Println(dep.Name)
		dep.Report()
	}
}

func (a *Atlas) HandleError(err error) {

	if err == nil {
		return
	}
	logrus.Fatal(err.Error())
}
