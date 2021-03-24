package server

import (
	bytes2 "bytes"
	json2 "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/mongodb/gridfs"
	rules2 "github.com/maczikasz/go-runs/internal/mongodb/rules"
	runbooks2 "github.com/maczikasz/go-runs/internal/mongodb/runbooks"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/sessions"
	"github.com/maczikasz/go-runs/internal/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func DoTestHttpRunbookManagementWithMongo(t *testing.T, client *mongodb.MongoClient) error {

	router, err := initializeRouter(client)
	if err != nil {
		return err
	}

	var step1Id string
	var step2Id string

	step1 := stepDTO{
		Summary: "Test summary",
		Markdown: MarkdownInfo{
			Content: "TEST_MARKDOWN",
			Type:    "gridfs",
		},
		Type: "Workaround",
	}

	step2 := stepDTO{
		Summary: "Test summary1",
		Markdown: MarkdownInfo{
			Content: "TEST_MARKDOWN1",
			Type:    "gridfs",
		},
		Type: "Workaround1",
	}

	Convey("When step1 is created", t, func() {

		bytes, _ := json2.Marshal(step1)

		r, _ := http.NewRequest(http.MethodPost, "/details", bytes2.NewReader(bytes))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		step1Id = w.Body.String()

		Convey("Then step 1 can be read back", func() {
			r, _ = http.NewRequest(http.MethodGet, "/details/"+step1Id, nil)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusOK)

			var details model.RunbookStepDetails

			body, _ := ioutil.ReadAll(w.Body)

			_ = json2.Unmarshal(body, &details)

			So(details.Markdown, ShouldEqual, step1.Markdown.Content)
			So(details.Type, ShouldEqual, step1.Type)
			So(details.Summary, ShouldEqual, step1.Summary)
			So(details.Id, ShouldEqual, step1Id)
		})
	})

	Convey("When step2 is created", t, func() {

		bytes, _ := json2.Marshal(step2)

		r, _ := http.NewRequest(http.MethodPost, "/details", bytes2.NewReader(bytes))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		step2Id = w.Body.String()

		Convey("Then steps can be read back", func() {

			r, _ = http.NewRequest(http.MethodGet, "/details/"+step2Id, nil)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusOK)

			body, _ := ioutil.ReadAll(w.Body)

			var details model.RunbookStepDetails

			_ = json2.Unmarshal(body, &details)

			So(details.Markdown, ShouldEqual, step2.Markdown.Content)
			So(details.Type, ShouldEqual, step2.Type)
			So(details.Summary, ShouldEqual, step2.Summary)
			So(details.Id, ShouldEqual, step2Id)
		})
	})

	Convey("When new runbook is created with the steps", t, func() {
		bytes, _ := json2.Marshal(runbookDTO{Steps: []string{step1Id, step2Id}})

		r, _ := http.NewRequest(http.MethodPost, "/runbooks", bytes2.NewReader(bytes))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		runbookId := w.Body.String()

		Convey("Then runbook steps could be read back", func() {

			r, _ := http.NewRequest(http.MethodGet, "/runbooks/"+runbookId, bytes2.NewReader(bytes))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusOK)

			var details model.RunbookDetails
			body, _ := ioutil.ReadAll(w.Body)

			_ = json2.Unmarshal(body, &details)

			So(details.Steps, test_utils.ShouldMatch, func(value interface{}) bool {
				if step, ok := value.(model.RunbookStepData); ok {
					return step.Id == step1Id && step.Summary == step1.Summary && step.Type == step1.Type
				}

				return false
			})

			So(details.Steps, test_utils.ShouldMatch, func(value interface{}) bool {
				if step, ok := value.(model.RunbookStepData); ok {
					return step.Id == step2Id && step.Summary == step2.Summary && step.Type == step2.Type
				}

				return false
			})
		})

	})

	return nil
}
func TestHttpRunbookManagementWithMongo(t *testing.T) {
	test_utils.RunMongoDBDockerTest(DoTestHttpRunbookManagementWithMongo, t)
}

func TestHttpServerWithMongo(t *testing.T) {
	test_utils.RunMongoDBDockerTest(DoTestHttpServerWithMongo, t)
}

func DoTestHttpServerWithMongo(t *testing.T, client *mongodb.MongoClient) error {
	router, err := initializeRouter(client)
	if err != nil {
		return err
	}

	Convey("Given server set up using mongodb", t, func() {

		Convey("When matching message rule is saved", func() {
			messageMatchingRunbookId := primitive.NewObjectID().Hex()
			rule := RuleCreateDTO{
				RuleType:    "message",
				MatcherType: "equal",
				RuleContent: "Test message",
				RunbookId:   messageMatchingRunbookId,
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
						So(session.Runbook.Id, ShouldEqual, messageMatchingRunbookId)
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

				nameMatchingRunbookId := primitive.NewObjectID().Hex()
				rule := RuleCreateDTO{
					RuleType:    "name",
					MatcherType: "contains",
					RuleContent: "name",
					RunbookId:   nameMatchingRunbookId,
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
							So(session.Runbook.Id, ShouldEqual, nameMatchingRunbookId)
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
							So(session.Runbook.Id, ShouldEqual, messageMatchingRunbookId)
						})
					})
				})
			})
		})
	})

	return nil
}

func initializeRouter(client *mongodb.MongoClient) (*gin.Engine, error) {
	config, err := rules2.LoadPriorityRuleConfigFromMongodb(client)
	//config:= initMatchers()
	//
	if err != nil {
		log.Fatalf("Failed to load rule config from mongodb: %s", err)
		return nil, err
	}

	runbookDataManager := runbooks2.RunbookDataManager{Client: client}
	runbookStepsDataManager := runbooks2.RunbookStepsDataManager{Client: client}
	fsClient, err := client.NewGridFSClient()

	if err != nil {
		log.Fatalf("Failed to load rule config from mongodb: %s", err)
		return nil, err
	}

	resolver := runbooks.MapRunbookMarkdownResolver{
		Resolvers: map[string]runbooks.MarkdownHandlers{"gridfs": {
			Resolver: gridfs.MarkdownResolver{Client: &gridfs.Client{Bucket: fsClient}},
			Writer:   gridfs.MarkdownWriter{Client: &gridfs.Client{Bucket: fsClient}}},
		},
	}

	runbookStepDetailsFinder := runbooks.RunbookStepDetailsFinder{
		RunbookStepsEntityFinder:    runbookStepsDataManager,
		RunbookStepMarkdownResolver: resolver,
	}

	stepDetailsWriter := runbooks.RunbookStepDetailsWriter{
		RunbookStepsEntityWriter:  runbookStepsDataManager,
		RunbookStepMarkdownWriter: resolver,
	}
	ruleManager := rules.FromMatcherConfig(config)
	sessionManager := sessions.NewInMemorySessionManager()
	runbookManager := runbooks.RunbookManager{RuleManager: ruleManager, RunbookFinder: runbookDataManager}
	errorManager := errors.DefaultErrorManager{
		SessionCreator: sessionManager,
		RunbookFinder:  runbookManager,
	}
	startupContext := StartupContext{
		RunbookDetailsFinder:     runbookDataManager,
		SessionStore:             sessionManager,
		RunbookStepDetailsFinder: runbookStepDetailsFinder,
		SessionFromErrorCreator:  errorManager,
		RuleSaver:                rules2.PersistentRuleWriter{Mongo: client},
		RuleFinder:               rules2.PersistentRuleReader{Mongo: client},
		RuleMatcher:              ruleManager,
		RunbookStepDetailsWriter: stepDetailsWriter,
		RunbookDetailsWriter: runbookDataManager,
		RuleReloader: func() {

			config, err := rules2.LoadPriorityRuleConfigFromMongodb(client)

			if err != nil {
				log.Fatalf("Failed to load rule config from mongodb: %s", err)
				panic(err.Error())
			}

			ruleManager.ReloadFromMatcherConfig(config)
		},
	}

	router := setupRouter(&startupContext, []string{})
	return router, nil
}
