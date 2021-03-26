package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/util"
)

//go:generate moq -out ../mocks/sessions.go --skip-ensure . SessionStore

type (
	SessionStore interface {
		GetSession(s string) (model.Session, error)
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
