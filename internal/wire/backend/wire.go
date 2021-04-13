//+build wireinject

package backend

import (
	"github.com/google/wire"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/wire/backend/config"
)

func InitializeStartupContext() (*server.StartupContext, func(), error) {
	panic(wire.Build(
		config.MongodbProviderSet,
		config.GridfsProviderSet,
		config.RunbookMarkdownResolverProviderSet,
		config.ErrorManagerProviderSet,
		config.RuleManagerProviderSet,
		config.RunbookStepDetailsFinderProviderSet,
		config.RunbookStepDetailsWriterProviderSet,
		config.RunbookStepsDataManagerProviderSet,
		config.RunbookDataManagerProviderSet,
		config.SessionManagerProviderSet,
		config.RunbookManagerProviderSet,
		config.PersistentRuleManagerProviderSet,
		config.ProvideRuleReloaderFunction,
		config.ProvideStartupContext,
	))
}
