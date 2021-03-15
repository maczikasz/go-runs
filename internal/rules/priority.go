package rules

import (
	"github.com/maczikasz/go-runs/internal/model"
)

type PriorityRuleManager struct {
	nameMatchers    []RuleRunbookPair
	messageMatchers []RuleRunbookPair
	tagMatchers     []RuleRunbookPair
}

const (
	REGEX    = "MatchAgainst"
	CONTAINS = "contains"
	EQUAL    = "equal"
)

type MatcherConfig struct {
	nameMatchers    map[string][]RuleRunbookPair
	messageMatchers map[string][]RuleRunbookPair
	tagMatchers     map[string][]RuleRunbookPair
}

func (r MatcherConfig) AddNameRegexMatchers(matchers map[string]RegexMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.nameMatchers[REGEX] = append(r.nameMatchers[REGEX], RuleRunbookPair{
			RunbookId: k,
			Rule:      NameRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddNameExactMatchers(matchers map[string]ExactMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.nameMatchers[EQUAL] = append(r.nameMatchers[EQUAL], RuleRunbookPair{
			RunbookId: k,
			Rule:      NameRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddNameContainsMatchers(matchers map[string]ContainsMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.nameMatchers[CONTAINS] = append(r.nameMatchers[CONTAINS], RuleRunbookPair{
			RunbookId: k,
			Rule:      NameRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddMessageRegexMatchers(matchers map[string]RegexMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.messageMatchers[REGEX] = append(r.messageMatchers[REGEX], RuleRunbookPair{
			RunbookId: k,
			Rule:      MessageRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddMessageExactMatchers(matchers map[string]ExactMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.messageMatchers[EQUAL] = append(r.messageMatchers[EQUAL], RuleRunbookPair{
			RunbookId: k,
			Rule:      MessageRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddMessageContainsMatchers(matchers map[string]ContainsMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.messageMatchers[CONTAINS] = append(r.messageMatchers[CONTAINS], RuleRunbookPair{
			RunbookId: k,
			Rule:      MessageRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddTagRegexMatchers(matchers map[string]RegexMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.tagMatchers[REGEX] = append(r.tagMatchers[REGEX], RuleRunbookPair{
			RunbookId: k,
			Rule:      TagRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddTagExactMatchers(matchers map[string]ExactMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.tagMatchers[EQUAL] = append(r.tagMatchers[EQUAL], RuleRunbookPair{
			RunbookId: k,
			Rule:      TagRule{innerMatcher: v},
		})
	}
	return &r
}

func (r MatcherConfig) AddTagContainsMatchers(matchers map[string]ContainsMatcher) *MatcherConfig {
	for k, v := range matchers {
		r.tagMatchers[CONTAINS] = append(r.tagMatchers[CONTAINS], RuleRunbookPair{
			RunbookId: k,
			Rule:      TagRule{innerMatcher: v},
		})
	}
	return &r
}

func NewMatcherConfig() *MatcherConfig {
	return &MatcherConfig{
		nameMatchers:    make(map[string][]RuleRunbookPair),
		messageMatchers: make(map[string][]RuleRunbookPair),
		tagMatchers:     make(map[string][]RuleRunbookPair),
	}
}

func FromMatcherConfig(c *MatcherConfig) *PriorityRuleManager {
	var orderedNameMatchers []RuleRunbookPair
	var orderedMessageMatchers []RuleRunbookPair
	var orderedTagMatchers []RuleRunbookPair

	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[EQUAL]...)
	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[CONTAINS]...)
	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[REGEX]...)

	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[EQUAL]...)
	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[CONTAINS]...)
	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[REGEX]...)

	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[EQUAL]...)
	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[CONTAINS]...)
	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[REGEX]...)

	return &PriorityRuleManager{
		nameMatchers:    orderedNameMatchers,
		messageMatchers: orderedMessageMatchers,
		tagMatchers:     orderedTagMatchers,
	}

}

func (p PriorityRuleManager) FindMatch(error2 model.Error) (string, bool) {
	matched, found := findMatcher(error2, p.nameMatchers)
	if found {
		return matched.RunbookId, true
	}

	matched, found = findMatcher(error2, p.messageMatchers)
	if found {
		return matched.RunbookId, true
	}

	matched, found = findMatcher(error2, p.tagMatchers)
	if found {
		return matched.RunbookId, true
	}

	return "", false
}

func findMatcher(error2 model.Error, matchers []RuleRunbookPair) (RuleRunbookPair, bool) {
	for _, matcher := range matchers {
		if matcher.Rule.Matches(error2) {
			return matcher, true
		}
	}

	return RuleRunbookPair{}, false
}
