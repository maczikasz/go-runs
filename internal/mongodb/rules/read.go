package rules

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	"go.mongodb.org/mongo-driver/bson"
)

type PersistentRuleReader struct {
	Mongo *mongodb.MongoClient
}

func (p PersistentRuleReader) FindOneRule(ruleType string, matcherType string, ruleContent string) (*model.RuleEntity, error) {
	collection, cancel, ctx := p.Mongo.Collection("rules")
	defer cancel()

	var result model.RuleEntity
	err := collection.FindOne(ctx, bson.D{{Key: "ruletype", Value: ruleType}, {Key: "matchertype", Value: matcherType}}).Decode(result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p PersistentRuleReader) ListRules() (*[]model.RuleEntity, error) {
	collection, cancel, ctx := p.Mongo.Collection("rules")
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var result []model.RuleEntity

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p PersistentRuleReader) ReadRegexMatchers(ruleType string, ruleCreator RegexRuleCreator) (*map[string]rules.RegexMatcher, error) {
	dtos, err := p.FindRule(ruleType, regex)
	if err != nil {
		return nil, err
	}

	result := make(map[string]rules.RegexMatcher)

	for _, dto := range *dtos {
		result[dto.RunbookId] = ruleCreator(dto.RuleContent)
	}

	return &result, nil
}

func (p PersistentRuleReader) ReadEqualsMatchers(ruleType string, ruleCreator ExactRuleCreator) (*map[string]rules.EqualsMatcher, error) {
	dtos, err := p.FindRule(ruleType, equal)
	if err != nil {
		return nil, err
	}

	result := make(map[string]rules.EqualsMatcher)

	for _, dto := range *dtos {
		result[dto.RunbookId] = ruleCreator(dto.RuleContent)
	}

	return &result, nil
}

func (p PersistentRuleReader) ReadContainsMatchers(ruleType string, ruleCreator ContainsRuleCreator) (*map[string]rules.ContainsMatcher, error) {
	dtos, err := p.FindRule(ruleType, contains)
	if err != nil {
		return nil, err
	}

	result := make(map[string]rules.ContainsMatcher)

	for _, dto := range *dtos {
		result[dto.RunbookId] = ruleCreator(dto.RuleContent)
	}

	return &result, nil
}

func (p PersistentRuleReader) FindRule(ruleType string, matcherType string) (*[]model.RuleEntity, error) {
	collection, cancel, ctx := p.Mongo.Collection("rules")
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{{Key: "ruletype", Value: ruleType}, {Key: "matchertype", Value: matcherType}})
	if err != nil {
		return nil, err
	}

	var result []model.RuleEntity

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
