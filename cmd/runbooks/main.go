package main

import (
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/sessions"
	log "github.com/sirupsen/logrus"
	"sync"
)

func initMatchers() *rules.MatcherConfig {
	config := rules.NewMatcherConfig()
	config = config.AddNameExactMatchers(map[string]rules.ExactMatcher{"test-1": {MatchAgainst: "Test Error 1"}})
	config = config.AddNameContainsMatchers(map[string]rules.ContainsMatcher{"test-1-2": {MatchAgainst: "Test Error"}})
	config = config.AddMessageContainsMatchers(map[string]rules.ContainsMatcher{"test-2": {MatchAgainst: "message"}})
	config = config.AddTagExactMatchers(map[string]rules.ExactMatcher{"test-3": {MatchAgainst: "match"}})
	return config
}
func main() {
	log.SetLevel(log.TraceLevel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	config := initMatchers()
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
	}

	server.StartHttpServer(&wg, &startupContext)
	wg.Wait()
}
