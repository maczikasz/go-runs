package server

import (
	"encoding/json"
	"github.com/maczikasz/go-runs/internal/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

//go:generate moq -out errors_test.go . SessionFromErrorCreator

//TODO shit name
type SessionFromErrorCreator interface {
	GetSessionForError(e model.Error) (string, error)
}

type incomingErrorHandler struct {
	errorHandler SessionFromErrorCreator
}

func (i incomingErrorHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodOptions {
		return
	}

	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Trace("Failed to read body of http request")
		http.Error(writer, "failed to read body", 500)
	}
	e := model.Error{}
	err = json.Unmarshal(bytes, &e)
	if err != nil {
		log.Tracef("Failed to parse json %s \n reason", err.Error())
		http.Error(writer, "failed to parse json", 500)
	}
	sessionId, err := i.errorHandler.GetSessionForError(e)

	if err != nil {
		log.Tracef("Failed to parse json %s \n reason", err.Error())
		//TOOD 404? 400?
		http.Error(writer, "failed to find proper runbook", 500)
	}

	_, err = writer.Write([]byte(sessionId))

	if err != nil {
		log.Warn("failed to write json to output")
	}

}
