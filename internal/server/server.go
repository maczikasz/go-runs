package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
}

func setupRouter(context *StartupContext, acceptedOrigins []string) *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowHeaders: []string{"Content-type", "Origin"},
		AllowMethods: []string{"POST", "GET"},
	}
	if len(acceptedOrigins) == 0 {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = acceptedOrigins
	}

	r.Use(cors.New(config))
	
	errorHandler := incomingErrorHandler{errorHandler: context.SessionFromErrorCreator}
	sessionHandler := sessionHandler{sessionStore: context.SessionStore}
	runbookHandler := runbookHandler{runbookDetailsFinder: context.RunbookDetailsFinder}
	runbookStepDetailsHandler := runbookStepDetailsHandler{context.RunbookStepDetailsFinder}
	r.POST("/errors", errorHandler.Serve)
	r.GET("/sessions/:sessionId", sessionHandler.Serve)
	r.GET("/runbooks/:runbookId", runbookHandler.Serve)
	r.GET("/details/:stepId", runbookStepDetailsHandler.Serve)

	return r
}
