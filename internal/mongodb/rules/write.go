package rules

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PersistentRuleWriter struct {
	Mongo *mongodb.MongoClient
}

//func (receiver PersistentRuleWriter) WriteRegexRule(ruleType string, matcher *rules.RegexMatcher, runbookId string) error {
//
//	return receiver.WriteRule(ruleType, regex, matcher.MatchAgainst.String(), runbookId)
//
//}
//
//func (receiver PersistentRuleWriter) WriteExactRule(ruleType string, matcher *rules.EqualsMatcher, runbookId string) error {
//
//	return receiver.WriteRule(ruleType, equal, matcher.MatchAgainst, runbookId)
//
//}
//
//func (receiver PersistentRuleWriter) WriteContainsRule(ruleType string, matcher *rules.ContainsMatcher, runbookId string) error {
//
//	return receiver.WriteRule(ruleType, contains, matcher.MatchAgainst, runbookId)
//
//}

func (receiver PersistentRuleWriter) WriteRule(ruleType string, matcherType string, ruleContent string, runbookId string) error {
	entity := model.RuleEntity{
		RuleType:    ruleType,
		MatcherType: matcherType,
		RuleContent: ruleContent,
		RunbookId:   runbookId,
	}

	collection, cancel, ctx := receiver.Mongo.Collection("rules")
	defer cancel()

	_, err := collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (receiver PersistentRuleWriter) DeleteRule(ruleId string) error {
	collection, cancel, ctx := receiver.Mongo.Collection("rules")
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(ruleId)
	if err != nil {
		return errors.Wrap(err, "invalid ID format for mongodb")
	}

	_, err = collection.DeleteOne(ctx, bson.D{{"_id", objectID}})

	if err != nil {
		return err
	}

	return nil
}

type RegexRuleCreator func(string) rules.RegexMatcher
type ExactRuleCreator func(string) rules.EqualsMatcher
type ContainsRuleCreator func(string) rules.ContainsMatcher
