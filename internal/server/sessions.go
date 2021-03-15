package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate moq -out sessions_test.go . SessionStore

type SessionStore interface {
	GetSession(s string) (model.Session, error)
}

type sessionHandler struct {
	sessionStore SessionStore
}

func (s sessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(request)
	session, err := s.sessionStore.GetSession(vars["sessionId"])

	err = util.HandleDataError(writer, request, err)
	if err != nil {
		return
	}

	err = util.WriteJsonResponse(writer, session)

	if err != nil {
		log.Warnf("Failed to find session %s", err)
	}

}
