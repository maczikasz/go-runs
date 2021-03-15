package test_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type data struct {
	field1 string
	field2 int
}

func TestShouldMatchMatcher(t *testing.T) {
	Convey("Given a list of items", t, func() {
		Convey("When the list contains a certain element matching a func", func() {
			items := []data{{
				field1: "a",
				field2: 1,
			}}
			Convey("Then the ShouldMatch matcher matches on string field", func() {
				So(items, ShouldMatch, func(value interface{}) bool {
					myData, ok := value.(data)
					if !ok {
						return false
					}
					return myData.field1 == "a"
				})
			})

			Convey("Then the ShouldMatch matcher matches on int field", func() {
				So(items, ShouldMatch, func(value interface{}) bool {
					myData, ok := value.(data)
					if !ok {
						return false
					}
					return myData.field2 == 1
				})
			})
		})
	})
}
