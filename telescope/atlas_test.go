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

func (suite *SuiteAtlas) SetupTest() {

	atlas, _ := NewAtlas("../go.mod", false, map[string]bool{}, map[string]OutdatedScope{}).(*Atlas)
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
