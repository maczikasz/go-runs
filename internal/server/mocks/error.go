// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	"github.com/maczikasz/go-runs/internal/model"
	"sync"
)

// ErrorManagerMock is a mock implementation of ErrorManager.
//
// 	func TestSomethingThatUsesErrorManager(t *testing.T) {
//
// 		// make and configure a mocked ErrorManager
// 		mockedErrorManager := &ErrorManagerMock{
// 			ManageErrorWitSessionFunc: func(e model.Error) (string, error) {
// 				panic("mock out the ManageErrorWitSession method")
// 			},
// 		}
//
// 		// use mockedErrorManager in code that requires ErrorManager
// 		// and then make assertions.
//
// 	}
type ErrorManagerMock struct {
	// ManageErrorWitSessionFunc mocks the ManageErrorWitSession method.
	ManageErrorWitSessionFunc func(e model.Error) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// ManageErrorWitSession holds details about calls to the ManageErrorWitSession method.
		ManageErrorWitSession []struct {
			// E is the e argument value.
			E model.Error
		}
	}
	lockManageErrorWitSession sync.RWMutex
}

// ManageErrorWitSession calls ManageErrorWitSessionFunc.
func (mock *ErrorManagerMock) ManageErrorWitSession(e model.Error) (string, error) {
	if mock.ManageErrorWitSessionFunc == nil {
		panic("ErrorManagerMock.ManageErrorWitSessionFunc: method is nil but ErrorManager.ManageErrorWitSession was just called")
	}
	callInfo := struct {
		E model.Error
	}{
		E: e,
	}
	mock.lockManageErrorWitSession.Lock()
	mock.calls.ManageErrorWitSession = append(mock.calls.ManageErrorWitSession, callInfo)
	mock.lockManageErrorWitSession.Unlock()
	return mock.ManageErrorWitSessionFunc(e)
}

// ManageErrorWitSessionCalls gets all the calls that were made to ManageErrorWitSession.
// Check the length with:
//     len(mockedErrorManager.ManageErrorWitSessionCalls())
func (mock *ErrorManagerMock) ManageErrorWitSessionCalls() []struct {
	E model.Error
} {
	var calls []struct {
		E model.Error
	}
	mock.lockManageErrorWitSession.RLock()
	calls = mock.calls.ManageErrorWitSession
	mock.lockManageErrorWitSession.RUnlock()
	return calls
}
