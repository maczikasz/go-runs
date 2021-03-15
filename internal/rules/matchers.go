package rules

import (
	"fmt"
	"regexp"
	"strings"
)

type stringMatcher interface {
	matches(value string) bool
}

type RegexMatcher struct {
	MatchAgainst *regexp.Regexp
}

func (r RegexMatcher) String() string {
	return fmt.Sprintf("regexp %s matches", r.MatchAgainst.String())
}

func (r RegexMatcher) matches(value string) bool {
	return r.MatchAgainst.MatchString(value)
}

type ContainsMatcher struct {
	MatchAgainst string
}

func (r ContainsMatcher) String() string {
	return fmt.Sprintf("string contains %s", r.MatchAgainst)
}

func (r ContainsMatcher) matches(value string) bool {
	return strings.Contains(value, r.MatchAgainst)
}

type ExactMatcher struct {
	MatchAgainst string
}

func (r ExactMatcher) String() string {
	return fmt.Sprintf("string is equal to %s", r.MatchAgainst)
}

func (r ExactMatcher) matches(value string) bool {
	return strings.EqualFold(value, r.MatchAgainst)
}
