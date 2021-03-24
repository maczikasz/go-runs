package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"sync"
)

func StartHttpServer(wg *sync.WaitGroup, context *StartupContext) {
	defer wg.Done()
	r := setupRouter(context, []string{"http://localhost:3000"})

	//TODO error handle
	_ = r.Run()
}

type StartupContext struct {
	RunbookDetailsFinder     RunbookDetailsFinder
	SessionStore             SessionStore
	RunbookStepDetailsFinder RunbookStepDetailsFinder
	SessionFromErrorCreator  SessionFromErrorCreator
	RuleSaver                RuleSaver
	RuleFinder               RuleFinder
	RuleMatcher              runbooks.RuleMatcher
	RunbookStepDetailsWriter RunbookStepWriter
	RunbookDetailsWriter     RunbookDetailsWriter
	RuleReloader             func()
}

func setupRouter(context *StartupContext, acceptedOrigins []string) *gin.Engine {
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

	errorHandler := incomingErrorHandler{errorHandler: context.SessionFromErrorCreator}
	sessionHandler := sessionHandler{sessionStore: context.SessionStore}
	runbookHandler := runbookHandler{
		runbookDetailsFinder:     context.RunbookDetailsFinder,
		runbookDetailsWriter:     context.RunbookDetailsWriter,
		runbookStepDetailsFinder: context.RunbookStepDetailsFinder,
	}
	runbookStepDetailsHandler := runbookStepDetailsHandler{
		runbookManager: context.RunbookStepDetailsFinder,
		stepWriter:     context.RunbookStepDetailsWriter,
	}
	ruleHandler := ruleHandler{
		ruleSaver:    context.RuleSaver,
		ruleFinder:   context.RuleFinder,
		ruleMatcher:  context.RuleMatcher,
		ruleReloader: context.RuleReloader,
	}
	r.POST("/rules", ruleHandler.AddNewRule)
	r.GET("/rules", ruleHandler.ListAllRules)
	r.DELETE("/rules/:ruleId", ruleHandler.DisableRule)
	r.PUT("/rules/:ruleId", ruleHandler.UpdateRule)
	r.GET("/rules/match", ruleHandler.TestRuleMatch)
	r.POST("/errors", errorHandler.SubmitError)
	r.GET("/sessions/:sessionId", sessionHandler.LookupSession)
	r.GET("/runbooks/:runbookId", runbookHandler.RetrieveRunbook)
	r.GET("/details/:stepId", runbookStepDetailsHandler.RetrieveRunbookStepDetails)
	r.POST("/details", runbookStepDetailsHandler.CreateNewStep)
	r.POST("/runbooks", runbookHandler.CreateNewRunbook)

	return r
}
