package main

import (
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	mongodb2 "github.com/maczikasz/go-runs/internal/rules/mongodb"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/sessions"
	log "github.com/sirupsen/logrus"
	"sync"
)

func initMatchers() *rules.PriorityMatcherConfig {
	config := rules.NewMatcherConfig()
	config = config.AddNameEqualsMatchers(&map[string]rules.EqualsMatcher{"test-1": {MatchAgainst: "Test Error 1"}})
	config = config.AddNameContainsMatchers(&map[string]rules.ContainsMatcher{"test-1-2": {MatchAgainst: "Test Error"}})
	config = config.AddMessageContainsMatchers(&map[string]rules.ContainsMatcher{"test-2": {MatchAgainst: "message"}})
	config = config.AddTagEqualsMatchers(&map[string]rules.EqualsMatcher{"test-3": {MatchAgainst: "match"}})
	return config
}

func main() {
	log.SetLevel(log.TraceLevel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	client, disconnectFunction := mongodb.InitializeMongoClient("mongodb://localhost:27017", "runs")
	defer disconnectFunction()

	config, err := mongodb2.LoadPriorityRuleConfigFromMongodb(client)
	//config:= initMatchers()
	//
	if err != nil {
		log.Fatalf("Failed to load rule config from mongodb: %s", err)
		panic(err.Error())
	}

	ruleManager := rules.FromMatcherConfig(config)
	sessionManager := sessions.NewInMemorySessionManager()
	runbookManager := runbooks.FakeRunbookManager{RuleManager: ruleManager}
	errorManager := errors.DefaultErrorManager{
		SessionCreator: sessionManager,
		RunbookFinder:  runbookManager,
	}
	startupContext := server.StartupContext{
		RunbookDetailsFinder:     runbookManager,
		SessionStore:             sessionManager,
		RunbookStepDetailsFinder: runbookManager,
		SessionFromErrorCreator:  errorManager,
		RuleSaver:                mongodb2.PersistentRuleWriter{Mongo: client},
		RuleFinder:               mongodb2.PersistentRuleReader{Mongo: client},
		RuleMatcher:              ruleManager,
		RuleReloader: func() {

			config, err := mongodb2.LoadPriorityRuleConfigFromMongodb(client)

			if err != nil {
				log.Fatalf("Failed to load rule config from mongodb: %s", err)
				panic(err.Error())
			}

			ruleManager.ReloadFromMatcherConfig(config)
		},
	}

	server.StartHttpServer(&wg, &startupContext)
	wg.Wait()
}
