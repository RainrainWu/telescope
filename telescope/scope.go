package telescope

import (
	"errors"
	"fmt"
	"strings"
)

type OutdatedScope int

const (
	UP_TO_DATE OutdatedScope = iota
	MAJOR
	MINOR
	PATCH
	UNKNOWN
)

var OutdatedScopeLiteral [5]string = [...]string{"UP_TO_DATE", "MAJOR", "MINOR", "PATCH", "UNKNOWN"}

var MapScopeColor map[OutdatedScope]int = map[OutdatedScope]int{
	UP_TO_DATE: 92,
	MAJOR:      91,
	MINOR:      93,
	PATCH:      94,
	UNKNOWN:    97,
}

func (o OutdatedScope) String() string {
	return OutdatedScopeLiteral[o]
}

func OutdatedScopeStrToEnum(scopeStr string) OutdatedScope {

	scopeStr = strings.ToUpper(scopeStr)
	for idx, scp := range OutdatedScopeLiteral {
		if scopeStr == scp {
			return OutdatedScope(idx)
		}
	}
	panic(fmt.Errorf("unknown scope %s", scopeStr))
}

func GetTopScope(scopes []OutdatedScope) OutdatedScope {

	if len(scopes) == 0 {
		panic(errors.New("should provide at least one item"))
	}

	var hold OutdatedScope = UNKNOWN
	for _, scp := range scopes {
		if scp < hold {
			hold = scp
		}
	}
	return hold
}
