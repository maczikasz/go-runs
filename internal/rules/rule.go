package rules

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/model"
)

// A Rule determines wheather an instance of an Error matches the given criteria
// represented by the Rule
type Rule interface {
	Matches(error model.Error) bool
}

type stringMatcher interface {
	matches(value string) bool
}

type MessageRule struct {
	innerMatcher stringMatcher
}

func (receiver MessageRule) String() string {
	return fmt.Sprintf("matching on message with %s", receiver.innerMatcher)
}

func (receiver MessageRule) Matches(error model.Error) bool {
	return receiver.innerMatcher.matches(error.Message)
}

type NameRule struct {
	innerMatcher stringMatcher
}

func (receiver NameRule) String() string {
	return fmt.Sprintf("matching on name with %s", receiver.innerMatcher)
}

func (receiver NameRule) Matches(error model.Error) bool {
	return receiver.innerMatcher.matches(error.Name)
}

type TagRule struct {
	innerMatcher stringMatcher
}

func (receiver TagRule) String() string {
	return fmt.Sprintf("matching on any tag with %s", receiver.innerMatcher)
}

func (receiver TagRule) Matches(error model.Error) bool {

	for _, v := range error.Tags {
		if receiver.innerMatcher.matches(v) {
			return true
		}
	}

	return false
}
