package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	log "github.com/sirupsen/logrus"
)

//go:generate moq -out mocks/errors_mock.go -skip-ensure . SessionFromErrorCreator

//TODO shit name
type SessionFromErrorCreator interface {
	GetSessionForError(e model.Error) (string, error)
}

type incomingErrorHandler struct {
	errorHandler SessionFromErrorCreator
}

func (i incomingErrorHandler) Serve(context *gin.Context) {

	e := model.Error{}
	err := context.BindJSON(&e)
	if err != nil {
		log.Tracef("Failed to parse json %s \n reason", err.Error())
		context.Status(400)
		_ = context.Error(err)
		return
	}
	sessionId, err := i.errorHandler.GetSessionForError(e)

	if err != nil {
		//TOOD 404? 400?
		context.Status(400)
		_ = context.Error(err)
		return
	}

	context.String(200, sessionId)

}
