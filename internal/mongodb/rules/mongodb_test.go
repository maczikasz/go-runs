package rules

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestMongoDBRulesLoadedCorrectly(t *testing.T) {
	test_utils.RunMongoDBDockerTest(DoTestMongoDBRulesLoadedCorrectly, t)
}

func DoTestMongoDBRulesLoadedCorrectly(t *testing.T, mongoClient *mongodb.MongoClient) error {

	Convey("Given mongoDB is connected", t, func() {

		writer := PersistentRuleWriter{Mongo: mongoClient}
		reader := PersistentRuleReader{Mongo: mongoClient}

		ruleTypes := []string{"name", "message", "tag"}
		for _, ruleType := range ruleTypes {
			testAndWriteExactRules(writer, ruleType, "test", reader)
			testAndWriteContainsRules(writer, ruleType, "test", reader)
			testAndWriteRegexRules(writer, ruleType, "test", reader)
		}

	})
	return nil
}

func testAndWriteExactRules(writer PersistentRuleWriter, ruleType string, content string, reader PersistentRuleReader) {
	Convey(fmt.Sprintf("When %s %s matcher is written", equal, ruleType), func() {
		err := writer.WriteRule(ruleType, equal, content, "1")

		So(err, ShouldBeNil)

		Convey("Then rule read back properly", func() {
			matchers, err2 := reader.ReadEqualsMatchers(ruleType, func(s string) rules.EqualsMatcher {
				return rules.EqualsMatcher{MatchAgainst: s}
			})

			So(err2, ShouldBeNil)
			So(len(*matchers), ShouldEqual, 1)
			So(matchers, test_utils.ShouldMatch, func(value interface{}) bool {
				if rule, ok := value.(rules.EqualsMatcher); ok {
					return rule.MatchAgainst == content
				}

				return false
			})
		})
	})
}

func testAndWriteRegexRules(writer PersistentRuleWriter, ruleType string, content string, reader PersistentRuleReader) {
	Convey(fmt.Sprintf("When %s %s matcher is written", regex, ruleType), func() {
		err := writer.WriteRule(ruleType, regex, content, "1")

		So(err, ShouldBeNil)

		Convey("Then rule read back properly", func() {
			matchers, err2 := reader.ReadRegexMatchers(ruleType, func(s string) rules.RegexMatcher {
				return rules.RegexMatcher{MatchAgainst: regexp.MustCompile(s)}
			})

			So(err2, ShouldBeNil)
			So(len(*matchers), ShouldEqual, 1)
			So(matchers, test_utils.ShouldMatch, func(value interface{}) bool {
				if rule, ok := value.(rules.RegexMatcher); ok {
					return rule.MatchAgainst.String() == content
				}

				return false
			})
		})
	})
}

func testAndWriteContainsRules(writer PersistentRuleWriter, ruleType string, content string, reader PersistentRuleReader) {
	Convey(fmt.Sprintf("When %s %s matcher is written", contains, ruleType), func() {
		err := writer.WriteRule(ruleType, contains, content, "1")

		So(err, ShouldBeNil)

		Convey("Then rule read back properly", func() {
			matchers, err2 := reader.ReadContainsMatchers(ruleType, func(s string) rules.ContainsMatcher {
				return rules.ContainsMatcher{MatchAgainst: s}
			})

			So(err2, ShouldBeNil)
			So(len(*matchers), ShouldEqual, 1)
			So(matchers, test_utils.ShouldMatch, func(value interface{}) bool {
				if rule, ok := value.(rules.ContainsMatcher); ok {
					return rule.MatchAgainst == content
				}

				return false
			})
		})
	})
}
