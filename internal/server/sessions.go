package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/sessions"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type sessionHandler struct {
	sessionManager sessions.SessionManager
}

func (s sessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {


	if request.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(request)
	session, err := s.sessionManager.GetSession(vars["sessionId"])

	err = HandleDataError(writer, request, err)
	if err != nil {
		return
	}

	err = WriteJsonResponse(writer, session)

	if err != nil {
		log.Warnf("Failed to find session %s", err)
	}

}
