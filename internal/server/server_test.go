package server

import (
	bytes2 "bytes"
	"encoding/json"
	"github.com/maczikasz/go-runs/internal/model"
	server "github.com/maczikasz/go-runs/internal/server/mocks"
	"github.com/maczikasz/go-runs/internal/test_utils"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var sessionStore = make(map[string]model.Session)

func TestStartHttpServer(t *testing.T) {

	testContext := StartupContext{

		RunbookDetailsFinder: &server.RunbookDetailsFinderMock{
			FindRunbookDetailsByIdFunc: func(id string) (model.RunbookDetails, error) {
				if id == "test-1" {
					return model.RunbookDetails{Steps: []model.RunbookStepData{
						{
							Id:      "rbs1",
							Summary: "Test step 1",
							Type:    "Workaround",
						},
					}}, nil
				}
				return model.RunbookDetails{}, model.CreateDataNotFoundError("runbook_details", id)
			},
		},
		SessionStore: &server.SessionStoreMock{
			GetSessionFunc: func(s string) (model.Session, error) {
				session, found := sessionStore[s]

				if !found {
					return model.Session{}, model.CreateDataNotFoundError("session", s)
				}

				return session, nil
			},
		},
		RunbookStepDetailsFinder: &server.RunbookStepDetailsFinderMock{
			FindRunbookStepDetailsByIdFunc: func(id string) (model.RunbookStepDetails, error) {
				if id == "rbs1" {
					return model.RunbookStepDetails{
						RunbookStepData: model.RunbookStepData{
							Summary: "Test workaround 1",
							Type:    "Workaround",
						},
						Markdown: "Test MD 1",
					}, nil
				}
				return model.RunbookStepDetails{}, model.CreateDataNotFoundError("step_details", id)
			},
		},
		SessionFromErrorCreator: &server.SessionFromErrorCreatorMock{
			GetSessionForErrorFunc: func(e model.Error) (string, error) {
				if e.Name == "Test Error 1" {
					sessionStore["s1"] = model.Session{
						Runbook:         model.RunbookRef{Id: "test-1"},
						SessionId:       "s1",
						Stats:           model.SessionStatistics{CompletedSteps: map[string]time.Time{}},
						TriggeringError: e,
					}
					return "s1", nil
				}

				return "", errors.New("No match for error")
			},
		},
	}

	Convey("Given http router is setup with fake context", t, func() {
		router := setupRouter(&testContext, []string{})
		var sessionId = ""
		Convey("When exact name matching error sent in", func() {
			err := model.Error{
				Name:    "Test Error 1",
				Message: "Test Message",
				Tags:    []string{"Tag1", "Tag2"},
			}
			bytes, _ := json.Marshal(err)
			r, _ := http.NewRequest(http.MethodPost, "/errors", bytes2.NewReader(bytes))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			body, _ := ioutil.ReadAll(w.Body)
			Convey("Then a sessionId is returned", func() {
				So(w.Code, ShouldEqual, 200)
				So(body, ShouldNotBeEmpty)
			})

			sessionId = string(body)

			Convey("When the sessionId is queried", func() {

				r, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionId, bytes2.NewReader(bytes))
				w := httptest.NewRecorder()

				router.ServeHTTP(w, r)

				body, _ := ioutil.ReadAll(w.Body)

				var session model.Session
				_ = json.Unmarshal(body, &session)

				Convey("Then the correct runbook id is returned", func() {

					So(w.Code, ShouldEqual, 200)
					So(session.SessionId, ShouldEqual, sessionId)
					So(session.Runbook.Id, ShouldEqual, "test-1")
				})

				Convey("When the details of the runbook are queried", func() {
					r, _ := http.NewRequest(http.MethodGet, "/runbooks/"+session.Runbook.Id, bytes2.NewReader(bytes))
					w := httptest.NewRecorder()

					router.ServeHTTP(w, r)

					var details model.RunbookDetails
					body, _ := ioutil.ReadAll(w.Body)

					_ = json.Unmarshal(body, &details)

					Convey("Then the correct runbook details are returned", func() {

						So(w.Code, ShouldEqual, 200)
						So(details.Steps, test_utils.ShouldMatch, func(value interface{}) bool {
							if step, ok := value.(model.RunbookStepData); ok {
								return step.Id == "rbs1" && step.Type == "Workaround" && step.Summary == "Test step 1"
							}

							return false
						})
					})
				})
			})

		})
	})
}
