package sessions

import (
	"github.com/google/uuid"
	"github.com/maczikasz/go-runs/internal/model"
	"sync"
	"time"
)

type InMemorySessionManager struct {
	rwLock   sync.RWMutex
	sessions map[string]model.Session
}

func (s *InMemorySessionManager) CreateNewSession(r model.RunbookRef, err model.Error) string {
	sessionId := uuid.New().String()
	newSession := model.Session{
		Runbook:   r,
		SessionId: sessionId,
		Stats: model.SessionStatistics{
			CompletedSteps: map[string]time.Time{},
		},
		TriggeringError: err,
	}
	s.sessions[sessionId] = newSession

	return sessionId
}

func NewInMemorySessionManager() *InMemorySessionManager {
	return &InMemorySessionManager{rwLock: sync.RWMutex{}, sessions: map[string]model.Session{}}
}

func (s *InMemorySessionManager) GetAllSessions() (result []model.Session, err error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	for _, v := range s.sessions {
		result = append(result, v)
	}

	if result == nil {
		return []model.Session{}, nil
	}
	return
}

func (s *InMemorySessionManager) GetSession(sessionId string) (model.Session, error) {

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	res, ok := s.sessions[sessionId]

	if !ok {
		return model.Session{}, model.CreateDataNotFoundError("session", sessionId)
	}

	return res, nil
}

func (s *InMemorySessionManager) UpdateSession(session model.Session) error {

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	storedSession, ok := s.sessions[session.SessionId]
	if !ok {
		return model.CreateDataNotFoundError("session", session.SessionId)
	}

	for k, v := range session.Stats.CompletedSteps {
		if storedSession.Stats.CompletedSteps[k] != v {
			storedSession.Stats.CompletedSteps[k] = v
		}
	}

	s.sessions[session.SessionId] = storedSession

	return nil
}
