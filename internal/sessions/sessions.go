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

func (s *InMemorySessionManager) CompleteStepInSession(sessionId string, stepId string, now time.Time) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	res, ok := s.sessions[sessionId]

	if !ok {
		return model.CreateDataNotFoundError("session", sessionId)
	}

	res.Stats.CompletedSteps[stepId] = now

	return nil
}

func (s *InMemorySessionManager) CreateNewSession(runbook model.RunbookRef, err model.Error) (string, error) {
	sessionId := uuid.New().String()
	newSession := model.Session{
		Runbook:   runbook,
		SessionId: sessionId,
		Stats: model.SessionStatistics{
			CompletedSteps: map[string]time.Time{},
		},
		TriggeringError: err,
	}
	s.sessions[sessionId] = newSession

	return sessionId, nil
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
