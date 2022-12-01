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
	"golang.org/x/mod/module"
)

const (
	proxyUrlGoModule      = "https://proxy.golang.org/%s/@v/list"
	proxyUrlPythonPackage = "https://pypi.org/pypi/%s/json"
)

type IDependable interface {
	QueryReleaseVersions(language Language, wg *sync.WaitGroup)
	GetOutdatedScope() OutdatedScope
}

type Dependency struct {
	Name                  string
	StrictSemVer          bool
	VersionCurrentLiteral string
	VersionCurrent        *semver.Version
	VersionLatest         *semver.Version
}

type PypiJson struct {
	Releases map[string][]struct{} `json:"releases"`
}

func NewSematicVersion(version string, strict bool) (*semver.Version, error) {

	semanticVersion, err := semver.NewVersion(version)
	if strict || err == nil {
		return semanticVersion, err
	}

	truncatedVersion := strings.FieldsFunc(
		version,
		func(r rune) bool {
			// truncate rc, alpha, and beta flags
			return strings.ContainsRune("rc a b", r)
		},
	)
	if len(truncatedVersion) == 0 {
		return nil, fmt.Errorf("invalid version string %s", version)
	}
	return semver.NewVersion(truncatedVersion[0])
}

func NewDependency(name, version string, strictSemVer bool) IDependable {

	versionCurrent, err := NewSematicVersion(version, strictSemVer)
	if err != nil {
		logrus.Debug(fmt.Sprintf("%s %s", err.Error(), version))
		return &Dependency{
			Name:                  name,
			StrictSemVer:          strictSemVer,
			VersionCurrentLiteral: version,
		}
	}

	return &Dependency{
		Name:                  name,
		StrictSemVer:          strictSemVer,
		VersionCurrentLiteral: version,
		VersionCurrent:        versionCurrent,
	}
}

func (d *Dependency) QueryReleaseVersions(language Language, wg *sync.WaitGroup) {

	defer wg.Done()

	if d.VersionCurrent == nil {
		return
	}

	switch language {
	case GO:
		d.queryVersionsGo()
	case PYTHON:
		d.queryVersionsPython()
	default:
		panic(fmt.Errorf("unsupported language %s", language.String()))
	}
}

func getVersionsResponse(url string) *http.Response {

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Fatal(err.Error())
		panic(fmt.Errorf("failed to build request with url %s", url))
	}
	request.Header.Set("User-Agent", "GoMajor/1.0")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Fatal(err.Error())
		panic(fmt.Errorf("failed to send request to url %s", url))
	}
	return response
}

func getLatestVersion(versions []string, strictSemVer bool) *semver.Version {

	versionsAvailable := semver.Collection{}
	for _, ver := range versions {
		newVersion, err := NewSematicVersion(strings.TrimSpace(ver), strictSemVer)
		if err != nil {
			logrus.Debug(fmt.Sprintf("invalid version %s", ver))
			continue
		}

		versionsAvailable = append(versionsAvailable, newVersion)
	}
	if len(versionsAvailable) == 0 {
		return nil
	}

	sort.Sort(versionsAvailable)
	return versionsAvailable[len(versionsAvailable)-1]
}

func (d *Dependency) queryVersionsGo() {

	modulePath, err := module.EscapePath(d.Name)
	if err != nil {
		logrus.Fatal(err.Error())
		panic(fmt.Errorf("failed to escape module path %s", d.Name))
	}
	response := getVersionsResponse(fmt.Sprintf(proxyUrlGoModule, modulePath))
	defer response.Body.Close()

	versionsBytes, _ := io.ReadAll(response.Body)
	versions := strings.Split(
		strings.TrimSpace(
			strings.ReplaceAll(string(versionsBytes), "\r\n", "\n"),
		),
		"\n",
	)
	d.VersionLatest = getLatestVersion(versions, d.StrictSemVer)
}

func (d *Dependency) queryVersionsPython() {

	response := getVersionsResponse(fmt.Sprintf(proxyUrlPythonPackage, d.Name))
	defer response.Body.Close()

	var pypiJson PypiJson
	body, _ := io.ReadAll(response.Body)
	err := json.Unmarshal(body, &pypiJson)
	if err != nil {
		return
	}

	versions := []string{}
	for ver := range pypiJson.Releases {
		versions = append(versions, ver)
	}
	d.VersionLatest = getLatestVersion(versions, d.StrictSemVer)
}

func (d *Dependency) GetOutdatedScope() OutdatedScope {

	current, latest := d.VersionCurrent, d.VersionLatest
	if current == nil || latest == nil {
		return UNKNOWN
	}
	if latest.Major() > current.Major() {
		return MAJOR
	}
	if latest.Minor() > current.Minor() {
		return MINOR
	}
	if latest.Patch() > current.Patch() {
		return PATCH
	}
	return UP_TO_DATE
}
