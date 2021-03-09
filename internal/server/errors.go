package server

import (
	"encoding/json"
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)


type incomingErrorHandler struct {
	errorManager errors.ErrorManager
}

func (i incomingErrorHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Trace("Failed to read body of http request")
		http.Error(res, "failed to read body", 500)
	}
	e := model.Error{}
	err = json.Unmarshal(bytes, &e)
	if err != nil {
		log.Tracef("Failed to parse json %s \n reason")
		http.Error(res, "failed to parse json", 500)
	}
	sessionId, err := i.errorManager.HandleIncomingError(e)

	if err != nil {
		log.Tracef("Failed to parse json %s \n reason")
		//TOOD 404? 400?
		http.Error(res, "failed to find proper runbook", 500)
	}

	_, err = res.Write([]byte(sessionId))

	if err != nil {
		log.Warn("failed to write json to output")
	}

}

