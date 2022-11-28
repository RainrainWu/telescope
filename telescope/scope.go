package telescope

import (
	"errors"
	"fmt"
	"strings"
)

type OutdatedScope int

const (
	MAJOR OutdatedScope = iota
	MINOR
	PATCH
	UP_TO_DATE
	UNKNOWN
)

func (o OutdatedScope) String() string {
	return [...]string{"MAJOR", "MINOR", "PATCH", "UP_TO_DATE", "UNKNOWN"}[o]
}

func OutdatedScopeStrToEnum(scopeStr string) OutdatedScope {

	scopeStr = strings.ToUpper(scopeStr)
	for idx, scp := range [...]string{"MAJOR", "MINOR", "PATCH"} {
		if scopeStr == scp {
			return OutdatedScope(idx)
		}
	}
	panic(errors.New(fmt.Sprintf("unknown scope %s", scopeStr)))
}
