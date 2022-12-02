package telescope

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
)

func TestNewSemanticVersion(t *testing.T) {

	params := []struct {
		name     string
		version  string
		strict   bool
		expected string
	}{
		{name: "classical", version: "v1.0.0", strict: true, expected: "v1.0.0"},
		{name: "incompatible", version: "65.0.0+incompatible", strict: true, expected: "65.0.0+incompatible"},
		{name: "date time", version: "v0.0.0-20170810143723-de5bf2ad4578", strict: true, expected: "v0.0.0-20170810143723-de5bf2ad4578"},
		{name: "rc tag", version: "v1.0.0rc0", strict: false, expected: "v1.0.0"},
		{name: "alpha tag", version: "v1.0.0a1", strict: false, expected: "v1.0.0"},
		{name: "beta tag", version: "v1.0.0b2", strict: false, expected: "v1.0.0"},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				ver_obtained, _ := NewSematicVersion(param.version, param.strict)
				ver_expected, _ := semver.NewVersion(param.expected)
				assert.Equal(t, ver_obtained, ver_expected)
			},
		)
	}
}

func TestNewSemanticVersionError(t *testing.T) {

	params := []struct {
		name    string
		version string
		strict  bool
	}{
		{name: "malformed", version: "v1..0.0", strict: false},
		{name: "strict rc tag", version: "v1.0.0rc0", strict: true},
		{name: "strict alpha tag", version: "v1.0.0a1", strict: true},
		{name: "strict beta tag", version: "v1.0.0b2", strict: true},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				ver_obtained, err := NewSematicVersion(param.version, param.strict)
				assert.Nil(t, ver_obtained)
				assert.NotNil(t, err)
			},
		)
	}
}

func TestGetLatestVersion(t *testing.T) {

	assert.Nil(t, getLatestVersion([]string{}, false))

	params := []struct {
		name     string
		versions []string
		expected string
	}{
		{name: "single versions", versions: []string{"v1.0.0"}, expected: "v1.0.0"},
		{name: "duplicated versions", versions: []string{"v1.0.0", "v1.0.0"}, expected: "v1.0.0"},
		{name: "different versions 1", versions: []string{"v1.2.0", "v1.0.0"}, expected: "v1.2.0"},
		{name: "different versions 2", versions: []string{"v1.1.0", "v2.0.0"}, expected: "v2.0.0"},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				ver_obtained := getLatestVersion(param.versions, true)
				ver_expected, _ := semver.NewVersion(param.expected)
				assert.Equal(t, ver_obtained, ver_expected)
			},
		)
	}
}

func TestGetOutdatedScope(t *testing.T) {

	params := []struct {
		name           string
		versionCurrent string
		versionLatest  string
		expected       OutdatedScope
	}{
		{name: "up to date", versionCurrent: "v1.0.0", versionLatest: "v1.0.0", expected: UP_TO_DATE},
		{name: "major", versionCurrent: "v1.0.0", versionLatest: "v2.0.0", expected: MAJOR},
		{name: "minor", versionCurrent: "v1.0.0", versionLatest: "v1.1.0", expected: MINOR},
		{name: "patch", versionCurrent: "v1.0.0", versionLatest: "v1.0.1", expected: PATCH},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				versionCurrent, _ := semver.NewVersion(param.versionCurrent)
				versionLatest, _ := semver.NewVersion(param.versionLatest)
				dep := Dependency{VersionCurrent: versionCurrent, VersionLatest: versionLatest}
				assert.Equal(t, dep.GetOutdatedScope(), param.expected)
			},
		)
	}
}
