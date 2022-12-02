package telescope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutdatedScopeString(t *testing.T) {

	params := []struct {
		name     string
		scope    OutdatedScope
		expected string
	}{
		{name: "up_to_date", scope: UP_TO_DATE, expected: "UP_TO_DATE"},
		{name: "major", scope: MAJOR, expected: "MAJOR"},
		{name: "minor", scope: MINOR, expected: "MINOR"},
		{name: "patch", scope: PATCH, expected: "PATCH"},
		{name: "unknown", scope: UNKNOWN, expected: "UNKNOWN"},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, param.scope.String(), param.expected)
			},
		)
	}
}

func TestOutdatedScopeStrToEnum(t *testing.T) {

	params := []struct {
		name     string
		scopeStr string
		expected OutdatedScope
	}{
		{name: "up_to_date", scopeStr: "up_to_date", expected: UP_TO_DATE},
		{name: "major", scopeStr: "major", expected: MAJOR},
		{name: "minor", scopeStr: "minor", expected: MINOR},
		{name: "patch", scopeStr: "patch", expected: PATCH},
		{name: "unknown", scopeStr: "unknown", expected: UNKNOWN},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, OutdatedScopeStrToEnum(param.scopeStr), param.expected)
			},
		)
	}
}

func TestGetTopScope(t *testing.T) {

	params := []struct {
		name     string
		scopes   []OutdatedScope
		expected OutdatedScope
	}{
		{name: "single scopes", scopes: []OutdatedScope{MAJOR}, expected: MAJOR},
		{name: "duplicated scopes", scopes: []OutdatedScope{MAJOR, MAJOR}, expected: MAJOR},
		{name: "different scopes 1", scopes: []OutdatedScope{MAJOR, MINOR}, expected: MAJOR},
		{name: "different scopes 2", scopes: []OutdatedScope{PATCH, MAJOR}, expected: MAJOR},
		{name: "different scopes 3", scopes: []OutdatedScope{PATCH, MINOR, PATCH}, expected: MINOR},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, GetTopScope(param.scopes), param.expected)
			},
		)
	}
}

func TestGetTopScopePanic(t *testing.T) {

	params := []struct {
		name   string
		scopes []OutdatedScope
	}{
		{name: "no scope", scopes: []OutdatedScope{}},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				assert.Panics(
					t,
					func() { GetTopScope(param.scopes) },
					"func did not panic",
				)
			},
		)
	}
}
