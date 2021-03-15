package server

import (
	bytes2 "bytes"
	"encoding/json"
	"github.com/maczikasz/go-runs/internal/context"
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/runbooks/test_data"
	"github.com/maczikasz/go-runs/internal/sessions"
	"github.com/maczikasz/go-runs/internal/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var config = initMatchers()

func initMatchers() *rules.MatcherConfig {
	config := rules.NewMatcherConfig()
	config = config.AddNameExactMatchers(map[string]rules.ExactMatcher{"test-1": {MatchAgainst: "Test Error 1"}})
	config = config.AddNameContainsMatchers(map[string]rules.ContainsMatcher{"test-1-2": {MatchAgainst: "Test Error"}})
	config = config.AddMessageContainsMatchers(map[string]rules.ContainsMatcher{"test-2": {MatchAgainst: "message"}})
	config = config.AddTagExactMatchers(map[string]rules.ExactMatcher{"test-3": {MatchAgainst: "match"}})
	return config
}

var sessionManager = sessions.NewInMemorySessionManager()
var runbookManager = test_data.FakeRunbookManager{RuleManager: rules.FromMatcherConfig(config)}
var errorManager = errors.DefaultErrorManager{}
var testContext = context.StartupContext{
	RuleManager:    rules.FromMatcherConfig(config),
	ErrorManager:   errors.WithManagers(&errorManager, sessionManager, runbookManager),
	SessionManager: sessionManager,
	RunbookManager: runbookManager,
}

func TestStartHttpServer(t *testing.T) {

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

				var details model.RunbookDetails

				Convey("When the details of the runbook are queried", func() {
					r, _ := http.NewRequest(http.MethodGet, "/runbooks/"+session.Runbook.Id, bytes2.NewReader(bytes))
					w := httptest.NewRecorder()

					router.ServeHTTP(w, r)

					body, _ := ioutil.ReadAll(w.Body)

					_ = json.Unmarshal(body, &details)

					Convey("Then the correct runbook details are returned", func() {

						So(w.Code, ShouldEqual, 200)
						So(details.Steps, test_utils.ShouldMatch, func(value interface{}) bool {
							if step, ok := value.(model.RunbookStepSummary); ok {
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
