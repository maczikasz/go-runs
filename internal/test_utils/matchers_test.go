package test_utils

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type data struct {
	field1 string
	field2 int
}

func TestShouldMatchListMatcher(t *testing.T) {
	Convey("Given a list of items", t, func() {
		Convey("When the list contains a certain element matching a func", func() {
			items := []data{{
				field1: "a",
				field2: 1,
			}}
			Convey("Then the ShouldMatch matcher matches on string field", func() {
				So(items, ShouldMatch, func(value interface{}) bool {
					fmt.Printf("%s\n", spew.Sdump(value))

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

func TestShouldMatchMapMatcher(t *testing.T) {
	Convey("Given a map of items", t, func() {
		Convey("When the map contains a certain element matching a func", func() {
			items := map[string]data{"a": {
				field1: "a",
				field2: 1,
			}}
			Convey("Then the ShouldMatch matcher matches on string field", func() {
				So(items, ShouldMatch, func(value interface{}) bool {
					myData, ok := value.(data)
					fmt.Printf("%s\n", spew.Sdump(value))
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

func TestShouldMatchListPtrMatcher(t *testing.T) {
	Convey("Given a list of items", t, func() {
		Convey("When the list contains a certain element matching a func", func() {
			items := []data{{
				field1: "a",
				field2: 1,
			}}
			Convey("Then the ShouldMatch matcher matches on string field", func() {
				So(&items, ShouldMatch, func(value interface{}) bool {
					myData, ok := value.(data)
					if !ok {
						return false
					}
					return myData.field1 == "a"
				})
			})

			Convey("Then the ShouldMatch matcher matches on int field", func() {
				So(&items, ShouldMatch, func(value interface{}) bool {
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

func TestShouldMatchMapPtrMapMatcher(t *testing.T) {
	Convey("Given a map of items", t, func() {
		Convey("When the map contains a certain element matching a func", func() {
			items := map[string]data{"a": {
				field1: "a",
				field2: 1,
			}}
			Convey("Then the ShouldMatch matcher matches on string field", func() {
				So(&items, ShouldMatch, func(value interface{}) bool {
					myData, ok := value.(data)
					if !ok {
						return false
					}
					return myData.field1 == "a"
				})
			})

			Convey("Then the ShouldMatch matcher matches on int field", func() {
				So(&items, ShouldMatch, func(value interface{}) bool {
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

func TestNilParameter(t *testing.T) {
	Convey("When should matcher supplied with nil value it fails", t, func() {
		result := ShouldMatch(nil, func(value interface{}) bool { return true })
		So(result, ShouldNotEqual, "")
	})
}
