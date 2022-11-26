package telescope

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/sirupsen/logrus"
)

type SemVerScope int

const (
	proxyUrlGoModule                  = "https://proxy.golang.org/%s/@v/list"
	proxyUrlPythonPackage             = "https://pypi.org/pypi/%s/json"
	MAJOR                 SemVerScope = iota
	MINOR
	PATCH
)

type Dependency struct {
	Name           string
	VersionCurrent *semver.Version
	VersionLatest  *semver.Version
}

type PypiJson struct {
	Releases map[string]struct{} `json:"releases"`
}

func NewDependency(name, version string, language Language) *Dependency {

	versionCurrent, err := semver.NewVersion(version)
	if err != nil {
		logrus.Fatal(err.Error(), version)
		panic(err)
	}

	return &Dependency{
		Name:           name,
		VersionCurrent: versionCurrent,
		VersionLatest:  nil,
	}
}

func (d *Dependency) QueryVersionsGo(wg *sync.WaitGroup) {

	defer wg.Done()

	url := fmt.Sprintf(proxyUrlGoModule, d.Name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "GoMajor/1.0")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	versionsBytes, _ := io.ReadAll(res.Body)
	versionsAvailable := semver.Collection{}
	versions := strings.Split(
		strings.TrimSpace(
			strings.ReplaceAll(string(versionsBytes), "\r\n", "\n"),
		),
		"\n",
	)
	for _, ver := range versions {
		newVersion, err := semver.NewVersion(strings.TrimSpace(ver))
		d.handleError(err)

		versionsAvailable = append(versionsAvailable, newVersion)
	}
	sort.Sort(versionsAvailable)
	d.VersionLatest = versionsAvailable[len(versionsAvailable)-1]
}

func (d *Dependency) QueryVersionsPython(wg *sync.WaitGroup) {

	defer wg.Done()

	url := fmt.Sprintf(proxyUrlPythonPackage, d.Name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "GoMajor/1.0")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var pypiJson PypiJson
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &pypiJson)

	versionsAvailable := semver.Collection{}
	for ver, _ := range pypiJson.Releases {
		newVersion, err := semver.NewVersion(strings.TrimSpace(ver))
		if err != nil {
			continue
		}

		versionsAvailable = append(versionsAvailable, newVersion)
	}
	sort.Sort(versionsAvailable)
	d.VersionLatest = versionsAvailable[len(versionsAvailable)-1]
}

func (d *Dependency) IsUpToDate(scope SemVerScope) bool {

	current, latest := d.VersionCurrent, d.VersionLatest
	switch scope {
	case MAJOR:
		if latest.Major() > current.Major() {
			return false
		}
	case MINOR:
		if latest.Minor() > current.Minor() {
			return false
		}
	case PATCH:
		if latest.Patch() > current.Patch() {
			return false
		}
	}
	return true
}

func (d *Dependency) Report() {

	ok := d.IsUpToDate(MINOR)
	fmt.Println(ok)
}

func (d *Dependency) handleError(err error) {

	if err == nil {
		return
	}
	logrus.Fatal(err.Error())
	panic(err)
}
