package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	log "github.com/sirupsen/logrus"
)

//go:generate moq -out ../mocks/error.go --skip-ensure . ErrorManager

//TODO shit name
type (
	ErrorManager interface {
		ManageErrorWitSession(e model.Error) (string, error)
	}

	IncomingErrorHandler struct {
		errorHandler ErrorManager
	}
)

func NewIncomingErrorHandler(creator ErrorManager) *IncomingErrorHandler {
	return &IncomingErrorHandler{errorHandler: creator}
}

func (i IncomingErrorHandler) SubmitError(context *gin.Context) {

	e := model.Error{}
	err := context.BindJSON(&e)
	if err != nil {
		log.Tracef("Failed to parse json %s \n reason", err.Error())
		context.Status(400)
		_ = context.Error(err)
		return
	}
	sessionId, err := i.errorHandler.ManageErrorWitSession(e)

	if err != nil {
		//TOOD 404? 400?
		context.Status(400)
		_ = context.Error(err)
		return
	}

	context.String(200, sessionId)

}
