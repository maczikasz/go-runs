// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"sync"
)

// RunbookStepDetailsFinderMock is a mock implementation of RunbookStepDetailsFinder.
//
// 	func TestSomethingThatUsesRunbookStepDetailsFinder(t *testing.T) {
//
// 		// make and configure a mocked RunbookStepDetailsFinder
// 		mockedRunbookStepDetailsFinder := &RunbookStepDetailsFinderMock{
// 			FindRunbookStepDetailsByIdFunc: func(id string) (model.RunbookStepDetails, error) {
// 				panic("mock out the FindRunbookStepDetailsById method")
// 			},
// 		}
//
// 		// use mockedRunbookStepDetailsFinder in code that requires RunbookStepDetailsFinder
// 		// and then make assertions.
//
// 	}
type RunbookStepDetailsFinderMock struct {
	// FindRunbookStepDetailsByIdFunc mocks the FindRunbookStepDetailsById method.
	FindRunbookStepDetailsByIdFunc func(id string) (model.RunbookStepDetails, error)

	// calls tracks calls to the methods.
	calls struct {
		// FindRunbookStepDetailsById holds details about calls to the FindRunbookStepDetailsById method.
		FindRunbookStepDetailsById []struct {
			// ID is the id argument value.
			ID string
		}
	}
	lockFindRunbookStepDetailsById sync.RWMutex
}

// FindRunbookStepDetailsById calls FindRunbookStepDetailsByIdFunc.
func (mock *RunbookStepDetailsFinderMock) FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error) {
	if mock.FindRunbookStepDetailsByIdFunc == nil {
		panic("RunbookStepDetailsFinderMock.FindRunbookStepDetailsByIdFunc: method is nil but RunbookStepDetailsFinder.FindRunbookStepDetailsById was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockFindRunbookStepDetailsById.Lock()
	mock.calls.FindRunbookStepDetailsById = append(mock.calls.FindRunbookStepDetailsById, callInfo)
	mock.lockFindRunbookStepDetailsById.Unlock()
	return mock.FindRunbookStepDetailsByIdFunc(id)
}

// FindRunbookStepDetailsByIdCalls gets all the calls that were made to FindRunbookStepDetailsById.
// Check the length with:
//     len(mockedRunbookStepDetailsFinder.FindRunbookStepDetailsByIdCalls())
func (mock *RunbookStepDetailsFinderMock) FindRunbookStepDetailsByIdCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockFindRunbookStepDetailsById.RLock()
	calls = mock.calls.FindRunbookStepDetailsById
	mock.lockFindRunbookStepDetailsById.RUnlock()
	return calls
}

// RunbookStepWriterMock is a mock implementation of RunbookStepWriter.
//
// 	func TestSomethingThatUsesRunbookStepWriter(t *testing.T) {
//
// 		// make and configure a mocked RunbookStepWriter
// 		mockedRunbookStepWriter := &RunbookStepWriterMock{
// 			WriteRunbookStepDetailsFunc: func(data model.RunbookStepData, markdown runbooks.Markdown, markdownLocationType string) (string, error) {
// 				panic("mock out the WriteRunbookStepDetails method")
// 			},
// 		}
//
// 		// use mockedRunbookStepWriter in code that requires RunbookStepWriter
// 		// and then make assertions.
//
// 	}
type RunbookStepWriterMock struct {
	// WriteRunbookStepDetailsFunc mocks the WriteRunbookStepDetails method.
	WriteRunbookStepDetailsFunc func(data model.RunbookStepData, markdown runbooks.Markdown, markdownLocationType string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// WriteRunbookStepDetails holds details about calls to the WriteRunbookStepDetails method.
		WriteRunbookStepDetails []struct {
			// Data is the data argument value.
			Data model.RunbookStepData
			// Markdown is the markdown argument value.
			Markdown runbooks.Markdown
			// MarkdownLocationType is the markdownLocationType argument value.
			MarkdownLocationType string
		}
	}
	lockWriteRunbookStepDetails sync.RWMutex
}

// WriteRunbookStepDetails calls WriteRunbookStepDetailsFunc.
func (mock *RunbookStepWriterMock) WriteRunbookStepDetails(data model.RunbookStepData, markdown runbooks.Markdown, markdownLocationType string) (string, error) {
	if mock.WriteRunbookStepDetailsFunc == nil {
		panic("RunbookStepWriterMock.WriteRunbookStepDetailsFunc: method is nil but RunbookStepWriter.WriteRunbookStepDetails was just called")
	}
	callInfo := struct {
		Data                 model.RunbookStepData
		Markdown             runbooks.Markdown
		MarkdownLocationType string
	}{
		Data:                 data,
		Markdown:             markdown,
		MarkdownLocationType: markdownLocationType,
	}
	mock.lockWriteRunbookStepDetails.Lock()
	mock.calls.WriteRunbookStepDetails = append(mock.calls.WriteRunbookStepDetails, callInfo)
	mock.lockWriteRunbookStepDetails.Unlock()
	return mock.WriteRunbookStepDetailsFunc(data, markdown, markdownLocationType)
}

// WriteRunbookStepDetailsCalls gets all the calls that were made to WriteRunbookStepDetails.
// Check the length with:
//     len(mockedRunbookStepWriter.WriteRunbookStepDetailsCalls())
func (mock *RunbookStepWriterMock) WriteRunbookStepDetailsCalls() []struct {
	Data                 model.RunbookStepData
	Markdown             runbooks.Markdown
	MarkdownLocationType string
} {
	var calls []struct {
		Data                 model.RunbookStepData
		Markdown             runbooks.Markdown
		MarkdownLocationType string
	}
	mock.lockWriteRunbookStepDetails.RLock()
	calls = mock.calls.WriteRunbookStepDetails
	mock.lockWriteRunbookStepDetails.RUnlock()
	return calls
}
