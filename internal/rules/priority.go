package rules

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/model"
	"sync"
)

type (
	RuleRunbookPair struct {
		RunbookId string
		Rule      Rule
	}

	PriorityRuleManager struct {
		ruleLock        *sync.RWMutex
		nameMatchers    []RuleRunbookPair
		messageMatchers []RuleRunbookPair
		tagMatchers     []RuleRunbookPair
	}

	PriorityMatcherConfig struct {
		nameMatchers    map[string][]RuleRunbookPair
		messageMatchers map[string][]RuleRunbookPair
		tagMatchers     map[string][]RuleRunbookPair
	}
)

const (
	regex    = "regex"
	contains = "contains"
	equal    = "equal"
)

func (m RuleRunbookPair) String() string {
	return fmt.Sprintf("rule: %s, runbook: %s", m.Rule, m.RunbookId)
}

func (r PriorityMatcherConfig) AddNameRegexMatchers(matchers *map[string]RegexMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.nameMatchers[regex] = append(r.nameMatchers[regex], RuleRunbookPair{
			RunbookId: k,
			Rule:      NameRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddNameEqualsMatchers(matchers *map[string]EqualsMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.nameMatchers[equal] = append(r.nameMatchers[equal], RuleRunbookPair{
			RunbookId: k,
			Rule:      NameRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddNameContainsMatchers(matchers *map[string]ContainsMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.nameMatchers[contains] = append(r.nameMatchers[contains], RuleRunbookPair{
			RunbookId: k,
			Rule:      NameRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddMessageRegexMatchers(matchers *map[string]RegexMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.messageMatchers[regex] = append(r.messageMatchers[regex], RuleRunbookPair{
			RunbookId: k,
			Rule:      MessageRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddMessageEqualsMatchers(matchers *map[string]EqualsMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.messageMatchers[equal] = append(r.messageMatchers[equal], RuleRunbookPair{
			RunbookId: k,
			Rule:      MessageRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddMessageContainsMatchers(matchers *map[string]ContainsMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.messageMatchers[contains] = append(r.messageMatchers[contains], RuleRunbookPair{
			RunbookId: k,
			Rule:      MessageRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddTagRegexMatchers(matchers *map[string]RegexMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.tagMatchers[regex] = append(r.tagMatchers[regex], RuleRunbookPair{
			RunbookId: k,
			Rule:      TagRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddTagEqualsMatchers(matchers *map[string]EqualsMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.tagMatchers[equal] = append(r.tagMatchers[equal], RuleRunbookPair{
			RunbookId: k,
			Rule:      TagRule{innerMatcher: v},
		})
	}
	return &r
}

func (r PriorityMatcherConfig) AddTagContainsMatchers(matchers *map[string]ContainsMatcher) *PriorityMatcherConfig {
	for k, v := range *matchers {
		r.tagMatchers[contains] = append(r.tagMatchers[contains], RuleRunbookPair{
			RunbookId: k,
			Rule:      TagRule{innerMatcher: v},
		})
	}
	return &r
}

func NewMatcherConfig() *PriorityMatcherConfig {
	return &PriorityMatcherConfig{
		nameMatchers:    make(map[string][]RuleRunbookPair),
		messageMatchers: make(map[string][]RuleRunbookPair),
		tagMatchers:     make(map[string][]RuleRunbookPair),
	}
}

func FromMatcherConfig(c *PriorityMatcherConfig) *PriorityRuleManager {
	var orderedNameMatchers []RuleRunbookPair
	var orderedMessageMatchers []RuleRunbookPair
	var orderedTagMatchers []RuleRunbookPair

	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[equal]...)
	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[contains]...)
	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[regex]...)

	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[equal]...)
	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[contains]...)
	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[regex]...)

	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[equal]...)
	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[contains]...)
	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[regex]...)

	return &PriorityRuleManager{
		ruleLock:        &sync.RWMutex{},
		nameMatchers:    orderedNameMatchers,
		messageMatchers: orderedMessageMatchers,
		tagMatchers:     orderedTagMatchers,
	}

}

func (r *PriorityRuleManager) ReloadFromMatcherConfig(c *PriorityMatcherConfig) {
	r.ruleLock.Lock()
	defer r.ruleLock.Unlock()

	var orderedNameMatchers []RuleRunbookPair
	var orderedMessageMatchers []RuleRunbookPair
	var orderedTagMatchers []RuleRunbookPair

	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[equal]...)
	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[contains]...)
	orderedNameMatchers = append(orderedNameMatchers, c.nameMatchers[regex]...)

	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[equal]...)
	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[contains]...)
	orderedMessageMatchers = append(orderedMessageMatchers, c.messageMatchers[regex]...)

	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[equal]...)
	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[contains]...)
	orderedTagMatchers = append(orderedTagMatchers, c.tagMatchers[regex]...)

	r.nameMatchers = orderedNameMatchers
	r.messageMatchers = orderedMessageMatchers
	r.tagMatchers = orderedTagMatchers

}

func (r *PriorityRuleManager) FindMatchingRunbook(error2 model.Error) (string, bool) {
	locker := r.ruleLock.RLocker()
	locker.Lock()
	defer locker.Unlock()

	matched, found := findMatcher(error2, r.nameMatchers)
	if found {
		return matched.RunbookId, true
	}

	matched, found = findMatcher(error2, r.messageMatchers)
	if found {
		return matched.RunbookId, true
	}

	matched, found = findMatcher(error2, r.tagMatchers)
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
