package server

import (
	bytes2 "bytes"
	json2 "encoding/json"
	"fmt"
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	mongodb2 "github.com/maczikasz/go-runs/internal/rules/mongodb"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/sessions"
	"github.com/maczikasz/go-runs/internal/test_utils"
	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpServerWithMongo(t *testing.T) {
	RunMongoDBDockerTest(DoTestHttpServerWithMongo, t)
}

func DoTestHttpServerWithMongo(t *testing.T, mongodburl string) error {
	client, disconnectFunction := mongodb.InitializeMongoClient(mongodburl, "local")
	defer disconnectFunction()

	config, err := mongodb2.LoadPriorityRuleConfigFromMongodb(client)
	//config:= initMatchers()
	//
	if err != nil {
		log.Fatalf("Failed to load rule config from mongodb: %s", err)
		return err
	}

	ruleManager := rules.FromMatcherConfig(config)
	sessionManager := sessions.NewInMemorySessionManager()
	//TODO user mongo runbook manager
	runbookManager := runbooks.FakeRunbookManager{RuleManager: ruleManager}
	errorManager := errors.DefaultErrorManager{
		SessionCreator: sessionManager,
		RunbookFinder:  runbookManager,
	}
	startupContext := StartupContext{
		RunbookDetailsFinder:     runbookManager,
		SessionStore:             sessionManager,
		RunbookStepDetailsFinder: runbookManager,
		SessionFromErrorCreator:  errorManager,
		RuleSaver:                mongodb2.PersistentRuleWriter{Mongo: client},
		RuleFinder:               mongodb2.PersistentRuleReader{Mongo: client},
		RuleMatcher:              ruleManager,
		RuleReloader: func() {

			config, err := mongodb2.LoadPriorityRuleConfigFromMongodb(client)

			if err != nil {
				log.Fatalf("Failed to load rule config from mongodb: %s", err)
				panic(err.Error())
			}

			ruleManager.ReloadFromMatcherConfig(config)
		},
	}

	router := setupRouter(&startupContext, []string{})

	Convey("Given server set up using mongodb", t, func() {
		collection, cancelFunc, ctx := client.Collection("rules")
		defer cancelFunc()
		_, _ = collection.DeleteMany(ctx, bson.D{})

		Convey("When matching message rule is saved", func() {
			rule := RuleCreateDTO{
				RuleType:    "message",
				MatcherType: "equal",
				RuleContent: "Test message",
				RunbookId:   "rb-1",
			}
			bytes, _ := json2.Marshal(&rule)

			r, _ := http.NewRequest(http.MethodPost, "/rules", bytes2.NewReader(bytes))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusCreated)

			Convey("Then rule list contains name rule", func() {
				r, _ := http.NewRequest(http.MethodGet, "/rules", bytes2.NewReader(bytes))
				w := httptest.NewRecorder()

				router.ServeHTTP(w, r)

				var rules []model.RuleEntity

				body, _ := ioutil.ReadAll(w.Body)

				_ = json2.Unmarshal(body, &rules)

				So(rules, ShouldNotBeEmpty)
				So(len(rules), ShouldEqual, 1)

				So(rules, test_utils.ShouldMatch, func(actual interface{}) bool {
					if matchRule, ok := actual.(model.RuleEntity); ok {
						matches := true

						matches = matches && matchRule.RunbookId == rule.RunbookId
						matches = matches && matchRule.RuleContent == rule.RuleContent
						matches = matches && matchRule.MatcherType == rule.MatcherType

						return matches
					}

					return false
				})

			})

			Convey("Then matching error returns session with correct runbookId", func() {

				err := model.Error{
					Name:    "Test Error 1",
					Message: "Test message",
					Tags:    []string{"Tag1", "Tag2"},
				}
				bytes, _ := json2.Marshal(err)
				r, _ := http.NewRequest(http.MethodPost, "/errors", bytes2.NewReader(bytes))
				w := httptest.NewRecorder()

				router.ServeHTTP(w, r)

				body, _ := ioutil.ReadAll(w.Body)
				Convey("Then a sessionId is returned", func() {
					So(w.Code, ShouldEqual, 200)
					So(body, ShouldNotBeEmpty)
				})

				sessionId := string(body)

				Convey("When the sessionId is queried", func() {

					r, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionId, bytes2.NewReader(bytes))
					w := httptest.NewRecorder()

					router.ServeHTTP(w, r)

					body, _ := ioutil.ReadAll(w.Body)

					var session model.Session
					_ = json2.Unmarshal(body, &session)

					Convey("Then the correct runbook id is returned", func() {

						So(w.Code, ShouldEqual, 200)
						So(session.SessionId, ShouldEqual, sessionId)
						So(session.Runbook.Id, ShouldEqual, "rb-1")
					})
				})
			})

			Convey("Then non-matching error returns 400", func() {
				err := model.Error{
					Name:    "Test Error 1",
					Message: "Different message",
					Tags:    []string{"Tag1", "Tag2"},
				}
				bytes, _ := json2.Marshal(err)
				r, _ := http.NewRequest(http.MethodPost, "/errors", bytes2.NewReader(bytes))
				w := httptest.NewRecorder()

				router.ServeHTTP(w, r)

				Convey("Then a sessionId is returned", func() {
					So(w.Code, ShouldEqual, 400)
				})
			})

			Convey("When matching name rule is saved", func() {

				rule := RuleCreateDTO{
					RuleType:    "name",
					MatcherType: "contains",
					RuleContent: "name",
					RunbookId:   "rb-0",
				}
				bytes, _ := json2.Marshal(&rule)

				r, _ := http.NewRequest(http.MethodPost, "/rules", bytes2.NewReader(bytes))
				w := httptest.NewRecorder()

				router.ServeHTTP(w, r)
				So(w.Code, ShouldEqual, http.StatusCreated)

				Convey("Then rule list contains both rules", func() {
					r, _ := http.NewRequest(http.MethodGet, "/rules", bytes2.NewReader(bytes))
					w := httptest.NewRecorder()

					router.ServeHTTP(w, r)

					var rules []model.RuleEntity

					body, _ := ioutil.ReadAll(w.Body)

					_ = json2.Unmarshal(body, &rules)

					So(rules, ShouldNotBeEmpty)
					So(len(rules), ShouldEqual, 2)
					So(rules, test_utils.ShouldMatch, func(actual interface{}) bool {
						if matchRule, ok := actual.(model.RuleEntity); ok {
							matches := true

							matches = matches && matchRule.RunbookId == rule.RunbookId
							matches = matches && matchRule.RuleContent == rule.RuleContent
							matches = matches && matchRule.MatcherType == rule.MatcherType

							return matches
						}

						return false
					})

				})

				Convey("Then matching error returns session with correct runbookId", func() {
					err := model.Error{
						Name:    "Test name 1",
						Message: "Test message",
						Tags:    []string{"Tag1", "Tag2"},
					}
					bytes, _ := json2.Marshal(err)
					r, _ := http.NewRequest(http.MethodPost, "/errors", bytes2.NewReader(bytes))
					w := httptest.NewRecorder()

					router.ServeHTTP(w, r)

					body, _ := ioutil.ReadAll(w.Body)
					Convey("Then a sessionId is returned", func() {
						So(w.Code, ShouldEqual, 200)
						So(body, ShouldNotBeEmpty)
					})

					sessionId := string(body)

					Convey("When the sessionId is queried", func() {

						r, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionId, bytes2.NewReader(bytes))
						w := httptest.NewRecorder()

						router.ServeHTTP(w, r)

						body, _ := ioutil.ReadAll(w.Body)

						var session model.Session
						_ = json2.Unmarshal(body, &session)

						Convey("Then the correct runbook id is returned", func() {

							So(w.Code, ShouldEqual, 200)
							So(session.SessionId, ShouldEqual, sessionId)
							So(session.Runbook.Id, ShouldEqual, "rb-0")
						})
					})
				})

				Convey("Then error matching the message rule returns with correct runbookId", func() {
					err := model.Error{
						Name:    "Test Error 1",
						Message: "Test message",
						Tags:    []string{"Tag1", "Tag2"},
					}
					bytes, _ := json2.Marshal(err)
					r, _ := http.NewRequest(http.MethodPost, "/errors", bytes2.NewReader(bytes))
					w := httptest.NewRecorder()

					router.ServeHTTP(w, r)

					body, _ := ioutil.ReadAll(w.Body)
					Convey("Then a sessionId is returned", func() {
						So(w.Code, ShouldEqual, 200)
						So(body, ShouldNotBeEmpty)
					})

					sessionId := string(body)

					Convey("When the sessionId is queried", func() {

						r, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionId, bytes2.NewReader(bytes))
						w := httptest.NewRecorder()

						router.ServeHTTP(w, r)

						body, _ := ioutil.ReadAll(w.Body)

						var session model.Session
						_ = json2.Unmarshal(body, &session)

						Convey("Then the correct runbook id is returned", func() {

							So(w.Code, ShouldEqual, 200)
							So(session.SessionId, ShouldEqual, sessionId)
							So(session.Runbook.Id, ShouldEqual, "rb-1")
						})
					})
				})
			})
		})
	})

	return nil
}

func RunMongoDBDockerTest(testFunction func(t *testing.T, mongoUrl string) error, t *testing.T) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	_ = resource.Expire(60)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		return testFunction(t, fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp")))
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}
