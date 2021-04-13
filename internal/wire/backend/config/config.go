package config

import (
	"github.com/google/wire"
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/mongodb"
	gridfs2 "github.com/maczikasz/go-runs/internal/mongodb/gridfs"
	rules2 "github.com/maczikasz/go-runs/internal/mongodb/rules"
	runbooks2 "github.com/maczikasz/go-runs/internal/mongodb/runbooks"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/server/handlers"
	mongodb2 "github.com/maczikasz/go-runs/internal/sessions/mongodb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MongoConfig struct {
	mongoUrl string
	database string
}

type GridfsConfig struct {
	gridfsMongoUrl string
	bucketName     string
}

type DataMongoDBClient mongodb.MongoClient
type GridfsMongoDBClient mongodb.MongoClient

func ProvideGridfsConfig() GridfsConfig {
	return GridfsConfig{
		gridfsMongoUrl: viper.GetString("markdown.gridfs.url"),
		bucketName:     viper.GetString("markdown.gridfs.bucket"),
	}
}
func ProvideGridfsMongoDBClient(config GridfsConfig) (*GridfsMongoDBClient, func(), error) {
	client, disconnectFunction, err := mongodb.InitializeMongoClient(config.gridfsMongoUrl, config.bucketName)

	return (*GridfsMongoDBClient)(client), disconnectFunction, err
}

var GridfsProviderSet = wire.NewSet(
	ProvideGridfsConfig,
	ProvideGridfsMongoDBClient,
)

func ProvideMongoConfig() MongoConfig {
	return MongoConfig{
		mongoUrl: viper.GetString("mongo.url"),
		database: viper.GetString("mongo.database"),
	}
}

func ProvideDataMongoDBClient(config MongoConfig) (*DataMongoDBClient, func(), error) {
	client, disconnectFunction, err := mongodb.InitializeMongoClient(config.mongoUrl, config.database)

	return (*DataMongoDBClient)(client), disconnectFunction, err
}

var MongodbProviderSet = wire.NewSet(
	ProvideMongoConfig,
	ProvideDataMongoDBClient,
)

func ProvidePriorityRuleMatcher(client *DataMongoDBClient) (*rules.PriorityRuleManager, error) {
	config, err := rules2.LoadPriorityRuleConfigFromMongodb((*mongodb.MongoClient)(client))
	if err != nil {
		return nil, err
	}

	ruleManager := rules.FromMatcherConfig(config)

	return ruleManager, nil
}

var RuleManagerProviderSet = wire.NewSet(
	ProvidePriorityRuleMatcher,
	wire.Bind(new(runbooks.RuleMatcher), new(*rules.PriorityRuleManager)),
)

func ProvideRunbookDataManager(client *DataMongoDBClient) *runbooks2.RunbookDataManager {
	return &runbooks2.RunbookDataManager{Client: (*mongodb.MongoClient)(client)}
}

var RunbookDataManagerProviderSet = wire.NewSet(
	ProvideRunbookDataManager,
	wire.Bind(new(runbooks.RunbookFinder), new(*runbooks2.RunbookDataManager)),
	wire.Bind(new(handlers.ReverseRunbookFinder), new(*runbooks2.RunbookDataManager)),
	wire.Bind(new(handlers.RunbookDetailsWriter), new(*runbooks2.RunbookDataManager)),
	wire.Bind(new(handlers.RunbookDetailsFinder), new(*runbooks2.RunbookDataManager)),
)

func ProvideRunbookStepsDataManager(client *DataMongoDBClient) *runbooks2.RunbookStepsDataManager {
	return &runbooks2.RunbookStepsDataManager{Client: (*mongodb.MongoClient)(client)}
}

var RunbookStepsDataManagerProviderSet = wire.NewSet(
	ProvideRunbookStepsDataManager,
	wire.Bind(new(runbooks.RunbookStepEntityWriter), new(*runbooks2.RunbookStepsDataManager)),
	wire.Bind(new(runbooks.RunbookStepEntityFinder), new(*runbooks2.RunbookStepsDataManager)),
)

type GridfsMarkdownHandlers runbooks.MarkdownHandlers

func ProvideGridFsClient(client *GridfsMongoDBClient) (*GridfsMarkdownHandlers, error) {
	fsClient, err := (*mongodb.MongoClient)(client).NewGridFSClient()

	if err != nil {
		return nil, err
	}

	return (*GridfsMarkdownHandlers)(&runbooks.MarkdownHandlers{
		Resolver: gridfs2.MarkdownResolver{Client: &gridfs2.Client{Bucket: fsClient}},
		Writer:   gridfs2.MarkdownWriter{Client: &gridfs2.Client{Bucket: fsClient}},
	}), nil
}

func ProvideRunbookMarkdownResolver(gridfsHandlers *GridfsMarkdownHandlers) *runbooks.MapRunbookMarkdownResolver {
	return runbooks.BuildNewMapRunbookMarkdownResolver(runbooks.Builder{{"gridfs", (*runbooks.MarkdownHandlers)(gridfsHandlers)}})
}

var RunbookMarkdownResolverProviderSet = wire.NewSet(
	ProvideGridFsClient,
	ProvideRunbookMarkdownResolver,
	wire.Bind(new(runbooks.RunbookStepMarkdownResolver), new(*runbooks.MapRunbookMarkdownResolver)),
	wire.Bind(new(runbooks.RunbookStepMarkdownWriter), new(*runbooks.MapRunbookMarkdownResolver)),
)

func ProvideStepDetailsFinder(runbookStepEntityFinder runbooks.RunbookStepEntityFinder, runbookStepMarkdownResolver runbooks.RunbookStepMarkdownResolver) *runbooks.RunbookStepDetailsFinder {
	return &runbooks.RunbookStepDetailsFinder{
		RunbookStepsEntityFinder:    runbookStepEntityFinder,
		RunbookStepMarkdownResolver: runbookStepMarkdownResolver,
	}
}

var RunbookStepDetailsFinderProviderSet = wire.NewSet(
	ProvideStepDetailsFinder,
	wire.Bind(new(handlers.RunbookStepDetailsFinder), new(*runbooks.RunbookStepDetailsFinder)),
)

func ProvideStepDetailsWriter(runbookStepsEntityWriter runbooks.RunbookStepEntityWriter, runbookStepMarkdownWriter runbooks.RunbookStepMarkdownWriter, runbookStepEntityFinder runbooks.RunbookStepEntityFinder) *runbooks.RunbookStepDetailsWriter {
	return &runbooks.RunbookStepDetailsWriter{
		RunbookStepsEntityWriter:  runbookStepsEntityWriter,
		RunbookStepMarkdownWriter: runbookStepMarkdownWriter,
		RunbookStepEntityFinder:   runbookStepEntityFinder,
	}
}

var RunbookStepDetailsWriterProviderSet = wire.NewSet(
	ProvideStepDetailsWriter,
	wire.Bind(new(handlers.RunbookStepWriter), new(*runbooks.RunbookStepDetailsWriter)),
)

func ProvideSessionManager(client *DataMongoDBClient) *mongodb2.SessionManager {
	return mongodb2.NewMongoDBSessionManager((*mongodb.MongoClient)(client))
}

var SessionManagerProviderSet = wire.NewSet(
	ProvideSessionManager,
	wire.Bind(new(errors.SessionCreator), new(*mongodb2.SessionManager)),
	wire.Bind(new(handlers.SessionStore), new(*mongodb2.SessionManager)),
)

func ProvideRunbookManager(ruleManager runbooks.RuleMatcher, runbookFinder runbooks.RunbookFinder) *runbooks.RunbookManager {
	return runbooks.NewRunbookManager(ruleManager, runbookFinder)
}

var RunbookManagerProviderSet = wire.NewSet(
	ProvideRunbookManager,
	wire.Bind(new(errors.RunbookFinder), new(*runbooks.RunbookManager)),
)

func ProvideDefaultErrorManager(sessionCreator errors.SessionCreator, runbookFinder errors.RunbookFinder) *errors.DefaultErrorManager {

	return errors.NewDefaultErrorManager(sessionCreator, runbookFinder)
}

var ErrorManagerProviderSet = wire.NewSet(
	ProvideDefaultErrorManager,
	wire.Bind(new(handlers.ErrorManager), new(*errors.DefaultErrorManager)),
)

func ProvideRuleReloaderFunction(manager *rules.PriorityRuleManager, client *DataMongoDBClient) func() {
	return func() {

		config, err := rules2.LoadPriorityRuleConfigFromMongodb((*mongodb.MongoClient)(client))

		if err != nil {
			log.Fatalf("Failed to load rule config from mongodb: %s", err)
			panic(err.Error())
		}

		manager.ReloadFromMatcherConfig(config)
	}
}

func ProvidePersistentRuleWriter(client *DataMongoDBClient) *rules2.PersistentRuleWriter {
	return &rules2.PersistentRuleWriter{Mongo: (*mongodb.MongoClient)(client)}
}

func ProvidePersistentRuleReader(client *DataMongoDBClient) *rules2.PersistentRuleReader {
	return &rules2.PersistentRuleReader{Mongo: (*mongodb.MongoClient)(client)}
}

var PersistentRuleManagerProviderSet = wire.NewSet(
	ProvidePersistentRuleWriter,
	wire.Bind(new(handlers.RuleSaver), new(*rules2.PersistentRuleWriter)),
	ProvidePersistentRuleReader,
	wire.Bind(new(handlers.RuleFinder), new(*rules2.PersistentRuleReader)),
)

func ProvideStartupContext(
	runbookDetailsFinder handlers.RunbookDetailsFinder,
	sessionStore handlers.SessionStore,
	runbookStepDetailsFinder handlers.RunbookStepDetailsFinder,
	errorManager handlers.ErrorManager,
	ruleSaver handlers.RuleSaver,
	ruleFinder handlers.RuleFinder,
	ruleMatcher runbooks.RuleMatcher,
	runbookStepDetailsWriter handlers.RunbookStepWriter,
	runbookDetailsWriter handlers.RunbookDetailsWriter,
	reverseRunbookFinder handlers.ReverseRunbookFinder,
	ruleReloader func(),
) *server.StartupContext {

	return &server.StartupContext{
		RunbookDetailsFinder:     runbookDetailsFinder,
		SessionStore:             sessionStore,
		RunbookStepDetailsFinder: runbookStepDetailsFinder,
		ErrorManager:             errorManager,
		RuleSaver:                ruleSaver,
		RuleFinder:               ruleFinder,
		RuleMatcher:              ruleMatcher,
		RunbookStepDetailsWriter: runbookStepDetailsWriter,
		RunbookDetailsWriter:     runbookDetailsWriter,
		ReverseRunbookFinder:     reverseRunbookFinder,
		RuleReloader:             ruleReloader,
	}
}
