package mongodb

import (
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	"regexp"
)

func LoadPriorityRuleConfigFromMongodb(client *mongodb.MongoClient) (*rules.PriorityMatcherConfig, error) {
	config := rules.NewMatcherConfig()

	reader := PersistentRuleReader{Mongo: client}

	containsMatchers, err := reader.ReadContainsMatchers("name", func(s string) rules.ContainsMatcher {
		return rules.ContainsMatcher{MatchAgainst: s}
	})

	if err != nil {
		return nil, err
	}

	config.AddNameContainsMatchers(containsMatchers)

	containsMatchers, err = reader.ReadContainsMatchers("message", func(s string) rules.ContainsMatcher {
		return rules.ContainsMatcher{MatchAgainst: s}
	})

	if err != nil {
		return nil, err
	}

	config.AddMessageContainsMatchers(containsMatchers)

	containsMatchers, err = reader.ReadContainsMatchers("tag", func(s string) rules.ContainsMatcher {
		return rules.ContainsMatcher{MatchAgainst: s}
	})

	if err != nil {
		return nil, err
	}

	config.AddTagContainsMatchers(containsMatchers)

	regexMatchers, err := reader.ReadRegexMatchers("name", func(s string) rules.RegexMatcher {
		//TODO sure that we want to panic?
		return rules.RegexMatcher{MatchAgainst: regexp.MustCompile(s)}
	})

	if err != nil {
		return nil, err
	}

	config.AddNameRegexMatchers(regexMatchers)

	regexMatchers, err = reader.ReadRegexMatchers("message", func(s string) rules.RegexMatcher {
		return rules.RegexMatcher{MatchAgainst: regexp.MustCompile(s)}
	})

	if err != nil {
		return nil, err
	}

	config.AddMessageRegexMatchers(regexMatchers)

	regexMatchers, err = reader.ReadRegexMatchers("tag", func(s string) rules.RegexMatcher {
		return rules.RegexMatcher{MatchAgainst: regexp.MustCompile(s)}
	})

	if err != nil {
		return nil, err
	}

	config.AddTagRegexMatchers(regexMatchers)

	equalsMatchers, err := reader.ReadEqualsMatchers("name", func(s string) rules.EqualsMatcher {
		return rules.EqualsMatcher{MatchAgainst: s}
	})

	if err != nil {
		return nil, err
	}

	config.AddNameEqualsMatchers(equalsMatchers)

	equalsMatchers, err = reader.ReadEqualsMatchers("message", func(s string) rules.EqualsMatcher {
		return rules.EqualsMatcher{MatchAgainst: s}
	})

	if err != nil {
		return nil, err
	}

	config.AddMessageEqualsMatchers(equalsMatchers)

	equalsMatchers, err = reader.ReadEqualsMatchers("tag", func(s string) rules.EqualsMatcher {
		return rules.EqualsMatcher{MatchAgainst: s}
	})

	if err != nil {
		return nil, err
	}

	config.AddTagEqualsMatchers(equalsMatchers)

	return config, nil
}
