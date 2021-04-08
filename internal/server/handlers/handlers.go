package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/util"
	"net/http"
	"time"
)

//go:generate moq -out ../mocks/sessions.go --skip-ensure . SessionStore

type (
	SessionStore interface {
		GetSession(s string) (model.Session, error)
		GetAllSessions() ([]model.Session, error)
		UpdateSession(session model.Session) error
		CompleteStepInSession(sessionId string, stepId string, now time.Time) error
	}

	SessionHandler struct {
		sessionStore SessionStore
	}
)

func NewSessionHandler(sessionStore SessionStore) *SessionHandler {
	return &SessionHandler{sessionStore: sessionStore}
}

func (s SessionHandler) LookupSession(context *gin.Context) {

	sessionId := context.Param("sessionId")
	session, err := s.sessionStore.GetSession(sessionId)

	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, session)
}

func (s SessionHandler) ListAllSessions(context *gin.Context) {
	sessions, err := s.sessionStore.GetAllSessions()
	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, sessions)

}

func (s SessionHandler) CompleteStepInSession(context *gin.Context) {
	sessionId := context.Param("sessionId")
	stepId := context.Param("stepId")

	err := s.sessionStore.CompleteStepInSession(sessionId, stepId, time.Now())

	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}
	context.Status(http.StatusOK)
}
