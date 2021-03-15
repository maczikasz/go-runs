package rules

import (
	"github.com/maczikasz/go-runs/internal/model"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestPriorityRuleManagerTypeOrder(t *testing.T) {
	e := model.Error{
		Name:    "name",
		Message: "message",
		Tags:    []string{"tag1", "tag2"},
	}

	Convey("Given a priority rule manager", t, func() {
		config := NewMatcherConfig()

		Convey("When all matchers are added", func() {
			config = config.AddNameContainsMatchers(map[string]ContainsMatcher{"test-1": {MatchAgainst: "name"}})
			config = config.AddMessageContainsMatchers(map[string]ContainsMatcher{"test-2": {MatchAgainst: "message"}})
			config = config.AddTagContainsMatchers(map[string]ContainsMatcher{"test-3": {MatchAgainst: "tag"}})

			Convey("Then Name matcher matches", func() {
				ruleManager := FromMatcherConfig(config)
				match, b := ruleManager.FindMatch(e)
				So(b, ShouldBeTrue)
				So(match, ShouldEqual, "test-1")
			})
		})

		Convey("When all message and tag matchers added", func() {
			config = config.AddMessageContainsMatchers(map[string]ContainsMatcher{"test-2": {MatchAgainst: "message"}})
			config = config.AddTagContainsMatchers(map[string]ContainsMatcher{"test-3": {MatchAgainst: "tag"}})

			Convey("Then Message matcher matches", func() {
				ruleManager := FromMatcherConfig(config)
				match, b := ruleManager.FindMatch(e)
				So(b, ShouldBeTrue)
				So(match, ShouldEqual, "test-2")
			})
		})

		Convey("When all tag matchers added", func() {
			config = config.AddTagContainsMatchers(map[string]ContainsMatcher{"test-3": {MatchAgainst: "tag"}})

			Convey("Then tag matcher matches", func() {
				ruleManager := FromMatcherConfig(config)
				match, b := ruleManager.FindMatch(e)
				So(b, ShouldBeTrue)
				So(match, ShouldEqual, "test-3")
			})
		})
	})

}

func TestPriorityRuleManagerMatcherOrderOrder(t *testing.T) {
	e := model.Error{
		Name:    "exact name",
		Message: "message",
		Tags:    []string{"tag1", "tag2"},
	}

	Convey("Given a priority rule manager", t, func() {

		config := NewMatcherConfig()
		Convey("When all type of matchers are added for Name", func() {
			config = config.AddNameExactMatchers(map[string]ExactMatcher{"test-1": {MatchAgainst: "exact name"}})
			config = config.AddNameContainsMatchers(map[string]ContainsMatcher{"test-2": {MatchAgainst: "name"}})
			config = config.AddNameRegexMatchers(map[string]RegexMatcher{"test-3": {MatchAgainst: regexp.MustCompile(".*")}})

			Convey("Then equals matcher matches", func() {
				ruleManager := FromMatcherConfig(config)
				match, b := ruleManager.FindMatch(e)
				So(b, ShouldBeTrue)
				So(match, ShouldEqual, "test-1")
			})
		})

		Convey("When contains and regex matchers are added for Name", func() {
			config = config.AddNameContainsMatchers(map[string]ContainsMatcher{"test-2": {MatchAgainst: "name"}})
			config = config.AddNameRegexMatchers(map[string]RegexMatcher{"test-3": {MatchAgainst: regexp.MustCompile(".*")}})

			Convey("Then contains matcher matches", func() {
				ruleManager := FromMatcherConfig(config)
				match, b := ruleManager.FindMatch(e)
				So(b, ShouldBeTrue)
				So(match, ShouldEqual, "test-2")
			})
		})

		Convey("When only regex matchers are added for Name", func() {
			config = config.AddNameRegexMatchers(map[string]RegexMatcher{"test-3": {MatchAgainst: regexp.MustCompile(".*")}})

			Convey("Then regex matcher matches", func() {
				ruleManager := FromMatcherConfig(config)
				match, b := ruleManager.FindMatch(e)
				So(b, ShouldBeTrue)
				So(match, ShouldEqual, "test-3")
			})
		})

	})

}
