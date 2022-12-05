package telescope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteAtlas struct {
	suite.Suite
	atlas *Atlas
}

func TestCompileRegExpRules(t *testing.T) {

	params := []struct{
		name string
		expressions []string
		expected int
	}{
		{name: "empty", expressions: []string{}, expected: 0},
		{name: "single expression", expressions: []string{"^github.com/.*$"}, expected: 1},
		{name: "multiple expressions", expressions: []string{"^github.com/.*$", "^k8s.io/.*$"}, expected: 2},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				patterns := compileRegExpRules(param.expressions)
				assert.Equal(t, len(patterns), param.expected)
			},
		)
	}
}

func TestMatchRegExpPatterns(t *testing.T) {

	patterns := compileRegExpRules([]string{"^github.com/.*$", "^k8s.io/.*$"})
	params := []struct{
		name string
		payload string
		expected bool
	}{
		{name: "not hit", payload: "gotest.tools/v3", expected: false},
		{name: "hit", payload: "k8s.io/api", expected: true},
	}
	for _, param := range params {

		t.Run(
			param.name,
			func(t *testing.T) {
				t.Parallel()
				patternHit := matchRegExpPatterns(patterns, param.payload)
				assert.Equal(t, patternHit, param.expected)
			},
		)
	}
}

func (suite *SuiteAtlas) SetupTest() {

	atlas, _ := NewAtlas("../go.mod", false, []string{}, map[OutdatedScope][]string{}).(*Atlas)
	suite.atlas = atlas
}

func (suite *SuiteAtlas) TestAppendDependency() {

	dep := NewDependency("module", "v1.0.0", false)
	suite.atlas.appendDependency(dep)
	assert.Contains(suite.T(), suite.atlas.dependencies, dep)
}

func (suite *SuiteAtlas) TestSortLexicographically() {

	dep_1 := NewDependency("_", "v1.0.0", false)
	suite.atlas.appendDependency(dep_1)
	suite.atlas.sortLexicographically()
	assert.Equal(suite.T(), suite.atlas.dependencies[0], dep_1)

	dep_2 := NewDependency("__", "v1.0.0", false)
	suite.atlas.appendDependency(dep_2)
	suite.atlas.sortLexicographically()
	assert.Equal(suite.T(), suite.atlas.dependencies[1], dep_2)
}

func TestSuiteAtlas(t *testing.T) {

	suite.Run(t, new(SuiteAtlas))
}
