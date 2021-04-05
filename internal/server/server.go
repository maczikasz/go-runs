package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server/handlers"
	"sync"
)

func StartHttpServer(wg *sync.WaitGroup, context *StartupContext) {
	defer wg.Done()
	r := SetupRouter(context, []string{"http://localhost:3000"})

	//TODO error handle
	_ = r.Run()
}

type StartupContext struct {
	RunbookDetailsFinder     handlers.RunbookDetailsFinder
	SessionStore             handlers.SessionStore
	RunbookStepDetailsFinder handlers.RunbookStepDetailsFinder
	ErrorManager             handlers.ErrorManager
	RuleSaver                handlers.RuleSaver
	RuleFinder               handlers.RuleFinder
	RuleMatcher              runbooks.RuleMatcher
	RunbookStepDetailsWriter handlers.RunbookStepWriter
	RunbookDetailsWriter     handlers.RunbookDetailsWriter
	ReverseRunbookFinder     handlers.ReverseRunbookFinder
	RuleReloader             func()
}

func SetupRouter(context *StartupContext, acceptedOrigins []string) *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowHeaders: []string{"Content-type", "Origin"},
		AllowMethods: []string{"POST", "GET", "DELETE", "PUT"},
	}
	if len(acceptedOrigins) == 0 {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = acceptedOrigins
	}

	r.Use(cors.New(config))

	errorHandler := handlers.NewIncomingErrorHandler(context.ErrorManager)
	sessionHandler := handlers.NewSessionHandler(context.SessionStore)
	runbookHandler := handlers.NewRunbookHandler(context.RunbookDetailsFinder, context.RunbookDetailsWriter, context.RunbookStepDetailsFinder)
	runbookStepDetailsHandler := handlers.NewRunbookStepDetailsHandler(context.RunbookStepDetailsFinder, context.RunbookStepDetailsWriter, context.ReverseRunbookFinder)
	ruleHandler := handlers.NewRuleHandler(context.RuleSaver, context.RuleFinder, context.RuleMatcher, context.RuleReloader)

	r.POST("/rules", ruleHandler.AddNewRule)
	r.GET("/rules", ruleHandler.ListAllRules)
	r.DELETE("/rules/:ruleId", ruleHandler.DisableRule)
	r.PUT("/rules/:ruleId", ruleHandler.UpdateRule)
	r.GET("/rules/match", ruleHandler.TestRuleMatch)

	r.POST("/errors", errorHandler.SubmitError)

	r.GET("/sessions/:sessionId", sessionHandler.LookupSession)
	r.PUT("/sessions/:sessionId/:stepId", sessionHandler.CompleteStepInSession)
	r.GET("/sessions", sessionHandler.ListAllSessions)

	r.GET("/details/:stepId", runbookStepDetailsHandler.RetrieveRunbookStepDetails)
	r.DELETE("/details/:stepId", runbookStepDetailsHandler.DeleteRunbookStep)
	r.PUT("/details/:stepId", runbookStepDetailsHandler.UpdateRunbookStep)
	r.POST("/details", runbookStepDetailsHandler.CreateNewStep)
	r.GET("/details", runbookStepDetailsHandler.ListAllSteps)

	r.POST("/runbooks", runbookHandler.CreateNewRunbook)
	r.GET("/runbooks", runbookHandler.ListAllRunbooks)
	r.GET("/runbooks/:runbookId", runbookHandler.RetrieveRunbook)
	r.DELETE("/runbooks/:runbookId", runbookHandler.DeleteRunbook)
	r.PUT("/runbooks/:runbookId", runbookHandler.UpdateRunbook)

	return r
}
