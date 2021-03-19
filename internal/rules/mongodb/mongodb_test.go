package mongodb

import (
	"context"
	"fmt"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/test_utils"
	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"regexp"
	"testing"
	"time"
)

func TestMongoDBRulesLoadedCorrectly(t *testing.T) {
	RunMongoDBDockerTest(DoTestMongoDBRulesLoadedCorrectly, t)
}

func DoTestMongoDBRulesLoadedCorrectly(t *testing.T, mongoUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return err
	}

	Convey("Given mongoDB is connected", t, func() {
		mongoClient, cancel := mongodb.InitializeMongoClient(mongoUrl, "local")
		defer cancel()
		writer := PersistentRuleWriter{Mongo: mongoClient}
		reader := PersistentRuleReader{Mongo: mongoClient}

		ruleTypes := []string{"name", "message", "tag"}
		for _, ruleType := range ruleTypes {
			testAndWriteExactRules(err, writer, ruleType, "test", reader)
			testAndWriteContainsRules(err, writer, ruleType, "test", reader)
			testAndWriteRegexRules(err, writer, ruleType, "test", reader)
		}

	})
	return nil
}

func testAndWriteExactRules(err error, writer PersistentRuleWriter, ruleType string, content string, reader PersistentRuleReader) {
	Convey(fmt.Sprintf("When %s %s matcher is written", equal, ruleType), func() {
		err = writer.WriteRule(ruleType, equal, content, "1")

		So(err, ShouldBeNil)

		Convey("Then rule read back properly", func() {
			matchers, err2 := reader.ReadEqualsMatchers(ruleType, func(s string) rules.EqualsMatcher {
				return rules.EqualsMatcher{MatchAgainst: s}
			})

			So(err2, ShouldBeNil)
			So(len(*matchers), ShouldEqual, 1)
			So(matchers, test_utils.ShouldMatch, func(value interface{}) bool {
				if rule, ok := value.(rules.EqualsMatcher); ok {
					return rule.MatchAgainst == content
				}

				return false
			})
		})
	})
}

func testAndWriteRegexRules(err error, writer PersistentRuleWriter, ruleType string, content string, reader PersistentRuleReader) {
	Convey(fmt.Sprintf("When %s %s matcher is written", regex, ruleType), func() {
		err = writer.WriteRule(ruleType, regex, content, "1")

		So(err, ShouldBeNil)

		Convey("Then rule read back properly", func() {
			matchers, err2 := reader.ReadRegexMatchers(ruleType, func(s string) rules.RegexMatcher {
				return rules.RegexMatcher{MatchAgainst: regexp.MustCompile(s)}
			})

			So(err2, ShouldBeNil)
			So(len(*matchers), ShouldEqual, 1)
			So(matchers, test_utils.ShouldMatch, func(value interface{}) bool {
				if rule, ok := value.(rules.RegexMatcher); ok {
					return rule.MatchAgainst.String() == content
				}

				return false
			})
		})
	})
}

func testAndWriteContainsRules(err error, writer PersistentRuleWriter, ruleType string, content string, reader PersistentRuleReader) {
	Convey(fmt.Sprintf("When %s %s matcher is written", contains, ruleType), func() {
		err = writer.WriteRule(ruleType, contains, content, "1")

		So(err, ShouldBeNil)

		Convey("Then rule read back properly", func() {
			matchers, err2 := reader.ReadContainsMatchers(ruleType, func(s string) rules.ContainsMatcher {
				return rules.ContainsMatcher{MatchAgainst: s}
			})

			So(err2, ShouldBeNil)
			So(len(*matchers), ShouldEqual, 1)
			So(matchers, test_utils.ShouldMatch, func(value interface{}) bool {
				if rule, ok := value.(rules.ContainsMatcher); ok {
					return rule.MatchAgainst == content
				}

				return false
			})
		})
	})
}

func DoTestMongoDBInDockerWorks(t *testing.T, mongoUrl string) error {
	var dockerLevelError error
	Convey("Given mongoDB is running inside a docker", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		Convey("When mongoDB connection is established", func() {
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

			if err != nil {
				log.Errorf("Failed to connect to mongodb %s", err.Error())
				dockerLevelError = err
			}

			Convey("Then mongoDB ping succeeds", func() {
				err = client.Ping(ctx, readpref.Primary())
				if err != nil {
					log.Errorf("Failed to ping mongodb %s", err.Error())

					dockerLevelError = err
				}
			})
		})

	})
	return dockerLevelError
}

func TestMongoDBInDockerWorks(t *testing.T) {
	RunMongoDBDockerTest(DoTestMongoDBInDockerWorks, t)
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
