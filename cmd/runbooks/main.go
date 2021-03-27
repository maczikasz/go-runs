package main

import (
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/mongodb/gridfs"
	rules2 "github.com/maczikasz/go-runs/internal/mongodb/rules"
	runbooks2 "github.com/maczikasz/go-runs/internal/mongodb/runbooks"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/sessions"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {
	log.SetLevel(log.TraceLevel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	client, disconnectFunction := mongodb.InitializeMongoClient("mongodb://localhost:27017", "runs")
	defer disconnectFunction()

	config, err := rules2.LoadPriorityRuleConfigFromMongodb(client)
	//config:= initMatchers()
	//
	if err != nil {
		log.Fatalf("Failed to load rule config from mongodb: %s", err)
		panic(err.Error())
	}

	runbookDataManager := runbooks2.RunbookDataManager{Client: client}
	runbookStepsDataManager := runbooks2.RunbookStepsDataManager{Client: client}
	fsClient, err := client.NewGridFSClient()

	if err != nil {
		log.Fatalf("Failed to load rule config from mongodb: %s", err)
		panic(err.Error())
	}

	resolver := runbooks.BuildNewMapRunbookMarkdownResolver(runbooks.Builder{
		{"gridfs", runbooks.MarkdownHandlers{
			Resolver: gridfs.MarkdownResolver{Client: &gridfs.Client{Bucket: fsClient}},
			Writer:   gridfs.MarkdownWriter{Client: &gridfs.Client{Bucket: fsClient}}},
		},
	})

	runbookStepDetailsFinder := runbooks.RunbookStepDetailsFinder{
		RunbookStepsEntityFinder:    runbookStepsDataManager,
		RunbookStepMarkdownResolver: resolver,
	}

	stepDetailsWriter := runbooks.RunbookStepDetailsWriter{
		RunbookStepsEntityWriter:  runbookStepsDataManager,
		RunbookStepMarkdownWriter: resolver,
	}
	ruleManager := rules.FromMatcherConfig(config)
	sessionManager := sessions.NewInMemorySessionManager()
	runbookManager := runbooks.NewRunbookManager(ruleManager, runbookDataManager)
	errorManager := errors.NewDefaultErrorManager(sessionManager, runbookManager)
	startupContext := server.StartupContext{
		RunbookDetailsFinder:     runbookDataManager,
		SessionStore:             sessionManager,
		RunbookStepDetailsFinder: runbookStepDetailsFinder,
		ErrorManager:             errorManager,
		RuleSaver:                rules2.PersistentRuleWriter{Mongo: client},
		RuleFinder:               rules2.PersistentRuleReader{Mongo: client},
		RuleMatcher:              ruleManager,
		RunbookStepDetailsWriter: stepDetailsWriter,
		RunbookDetailsWriter:     runbookDataManager,
		RuleReloader: func() {

			config, err := rules2.LoadPriorityRuleConfigFromMongodb(client)

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
